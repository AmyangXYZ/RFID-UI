package main

import (
	"fmt"
	"math"
	"time"
)

type Tag struct {
	ChDataFromReader chan RFIDData `json:"-"`
	ChSigBreak       chan bool     `json:"-"` // use - to omit it when convert to JSON
	EPC              string        `json:"epc"`
	EPC24            string        `json:"epc24"`
	Data             []RFIDData    `json:"data"`

	AddPortFlag bool    `json:"add_port_flag"`
	LED         string  `json:"led"`
	Dist        float64 `json:"dist"`
	Starttime   int64   `json:"starttime"`
	Endtime     int64   `json:"endtime"`
	MeanPower   int     `json:"meanpower"`
}

func newTag(epc string, epc24 string) *Tag {
	return &Tag{
		ChDataFromReader: make(chan RFIDData, 1024),
		ChSigBreak:       make(chan bool),
		EPC:              epc,
		EPC24:            epc24,
		AddPortFlag:      false,
		LED:              "GREY",
		Dist:             Distance,
	}
}

func (tag *Tag) countPortNumType() map[int]int { //return {1:11, 9:1} ==> len=2
	m := make(map[int]int) // dictionary is called map in golang
	for _, data := range tag.Data {
		m[data.AntennaPort] = m[data.AntennaPort] + 1
	}
	return m
}

func (tag *Tag) maxIndexPower() (int, int) {
	// powerList := make([]int, 0)
	maxValue, minValue, maxIndex := -99, 0, 0

	firstValue := true

	for i, data := range tag.Data {

		if firstValue {
			minValue = data.PeakRssi
			maxValue = data.PeakRssi
			firstValue = false
		} else {
			if data.PeakRssi < minValue {
				minValue = data.PeakRssi
			}
			if data.PeakRssi > maxValue {
				maxValue = data.PeakRssi
				maxIndex = i
			}
		}
	}
	meanPower := int(maxValue+minValue) / 2
	print("in function", maxValue, minValue, meanPower)
	return maxIndex, meanPower
}

func (tag *Tag) handleData() {

	for {
		select {
		case data := <-tag.ChDataFromReader:
			// print("initial len(tag.Data)", len(tag.Data))
			if (len(tag.Data) == 0) && (data.AntennaPort == 17) {
				tag.AddPortFlag = true
				fmt.Println("Antenna 17 active", tag.AddPortFlag)
				ChLEDToUI <- LEDServerToUI{tag.EPC, "GREEN"}
				tag.LED = "GREEN"
			}
			if data.AntennaPort == 17 && tag.AddPortFlag {
				tag.Data = append(tag.Data, data)
				// fmt.Println("==17 port adding", tag.Data)
			}
			if data.AntennaPort == 9 && tag.AddPortFlag {
				tag.Data = append(tag.Data, data)
				// fmt.Println("==9 port adding", tag.Data, len(tag.countPortNumType()))
				if len(tag.countPortNumType()) > 1 {
					maxIndex, meanPower := tag.maxIndexPower()
					tag.MeanPower = meanPower
					fmt.Println("==9 port adding before clean, should see 2 ports", tag.Data, len(tag.countPortNumType()))
					fmt.Println("meanPower", meanPower)

					tag.Dist = Distance

					tag.Starttime = tag.Data[maxIndex].FirstSeenTimestamp
					fmt.Println("tag.Starttime", tag.Starttime)
					tag.Data = []RFIDData{}

				}
				if data.PeakRssi > tag.MeanPower {
					timeEnd := data.FirstSeenTimestamp
					fmt.Println("timeEnd", timeEnd)
					timeDiff := float64(timeEnd-tag.Starttime) / 1000000
					speed := float64(Distance / timeDiff)
					speed = math.Round(speed*1000) / 1000
					currentTime := time.Now()
					timeText := currentTime.Format("15:04:05 PM")
					ChLEDToUI <- LEDServerToUI{tag.EPC, "RED"}
					tag.LED = "RED"
					ChDataToUI <- DataServerToUI{tag.EPC, speed, timeText, tag.Dist}
					fmt.Println("speed:", tag.EPC, speed, Distance, timeDiff, timeEnd)
					tag.Data = []RFIDData{}

					tag.AddPortFlag = false
					go func() {
						time.Sleep(time.Duration(timeDiff) * time.Second)
						fmt.Println(tag.EPC, "ready to run")
						ChLEDToUI <- LEDServerToUI{tag.EPC, "GREY"}
						tag.LED = "GREY"
					}()
				}

			}

			// 	// go func() {
			// 	// 	time.Sleep(time.Duration(2*timeDiff) * time.Second)
			// 	// 	tag.AddPortFlag = true
			// 	// 	fmt.Println(tag.EPC, "timer end, ready to run")
			// 	// 	ChLEDToUI <- LEDServerToUI{tag.EPC, "GREY"}
			// 	// 	tag.LED = "GREY"
			// 	// }()
			// }
		case <-tag.ChSigBreak:
			fmt.Println("quit", tag.EPC)
			return
		}
	}
}
