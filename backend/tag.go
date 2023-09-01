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
	AddPortFlag      bool          `json:"add_port_flag"`
	LED              string        `json:"led"`
}

func newTag(epc string, epc24 string) *Tag {
	return &Tag{
		ChDataFromReader: make(chan RFIDData, 1024),
		ChSigBreak:       make(chan bool),
		EPC:              epc,
		EPC24:            epc24,
		AddPortFlag:      true,
		LED:              "GREY",
	}
}

func (tag *Tag) countPortNumType() map[int]int { //return {1:11, 9:1} ==> len=2
	m := make(map[int]int) // dictionary is called map in golang
	for _, data := range tag.Data {
		m[data.AntennaPort] = m[data.AntennaPort] + 1
	}
	return m
}

func (tag *Tag) handleData() {
	for {
		select {
		case data := <-tag.ChDataFromReader:
			if tag.AddPortFlag {
				ChLEDToUI <- LEDServerToUI{tag.EPC, "GREEN"}
				tag.LED = "GREEN"
				tag.Data = append(tag.Data, data)
				// fmt.Println(tag.EPC, tag.Data)
			}
			if len(tag.countPortNumType()) > 1 {
				ChLEDToUI <- LEDServerToUI{tag.EPC, "RED"}
				tag.LED = "RED"

				tag.AddPortFlag = false
				timeRangeStart := tag.Data[0].FirstSeenTimestamp
				timeRangeEnd := tag.Data[len(tag.Data)-1].FirstSeenTimestamp
				timeDiff := float64(timeRangeEnd-timeRangeStart) / 1000000
				speed := float64(Distance / timeDiff)
				speed = math.Round(speed*1000) / 1000
				// hour, min, sec := time.Now().Clock()

				// timeText := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)
				currentTime := time.Now()
				timeText := currentTime.Format("15:04:05 PM")

				ChDataToUI <- DataServerToUI{tag.EPC, speed, timeText}
				fmt.Println("speed:", tag.EPC, speed)

				tag.Data = []RFIDData{}

				fmt.Println(tag.EPC, "timer start")
				go func() {
					time.Sleep(time.Duration(2*timeDiff) * time.Second)
					tag.AddPortFlag = true
					fmt.Println(tag.EPC, "timer end, ready to run")
					ChLEDToUI <- LEDServerToUI{tag.EPC, "GREY"}
					tag.LED = "GREY"

				}()

			}
		case <-tag.ChSigBreak:
			fmt.Println("quit", tag.EPC)

			return
		}
	}
}
