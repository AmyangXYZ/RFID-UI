package main

import (
	"fmt"
	"math"
	"time"
)

const WindowSize = 5
const PeakDropThreshold = 1

type RSSIEntry struct {
	PeakRssi           int
	FirstSeenTimestamp int64
	Port               int
}

type Tag struct {
	ChDataFromReader chan RFIDData `json:"-"`
	ChSigBreak       chan bool     `json:"-"` // use - to omit it when convert to JSON
	EPC              string        `json:"epc"`
	EPC24            string        `json:"epc24"`
	Data             []RFIDData    `json:"data"`
	RSSIWindow       []RSSIEntry
	AddPortFlag      bool    `json:"add_port_flag"`
	LED              string  `json:"led"`
	Dist             float64 `json:"dist"`
	Starttime        int64   `json:"starttime"`
	Endtime          int64   `json:"endtime"`
	MeanPower        int     `json:"meanpower"`
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
func (tag *Tag) reverseData() {
	for i, j := 0, len(tag.Data)-1; i < j; i, j = i+1, j-1 {
		tag.Data[i], tag.Data[j] = tag.Data[j], tag.Data[i]
	}
}

func (tag *Tag) maxIndexPower() {
	//use in first antenna. Have global array, moving window search in reverse order
	tag.reverseData()

	for _, item := range tag.Data {
		tag.updateWindow(item)
		if tag.detectFirstPeak() {
			tag.Starttime = item.FirstSeenTimestamp
			break
		}
	}
}

func (tag *Tag) updateWindow(newData RFIDData) {
	if len(tag.RSSIWindow) >= WindowSize {
		tag.RSSIWindow = tag.RSSIWindow[1:] // Remove the oldest data
	}
	tag.RSSIWindow = append(tag.RSSIWindow, RSSIEntry{
		PeakRssi:           newData.PeakRssi,
		FirstSeenTimestamp: newData.FirstSeenTimestamp,
		Port:               newData.AntennaPort,
	})
}

func (tag *Tag) detectFirstPeak() bool { // first peak for 1st antenna
	if len(tag.RSSIWindow) < WindowSize {
		return false // Not enough data
	}

	maxVal := tag.RSSIWindow[0].PeakRssi
	maxIndex := 0

	// Find max RSSI within the window
	for i, entry := range tag.RSSIWindow {
		if entry.PeakRssi > maxVal {
			maxVal = entry.PeakRssi
			maxIndex = i
		}
	}

	// Ensure a significant drop after the peak
	if maxIndex < len(tag.RSSIWindow)-1 {
		if maxVal-tag.RSSIWindow[maxIndex+1].PeakRssi >= PeakDropThreshold {
			return true
		}
	}

	return false
}
func (tag *Tag) detectLastPeak() bool { // last peak search in forward order for 2nd antenna
	if len(tag.RSSIWindow) < WindowSize {
		return false // Not enough data
	}

	maxVal := tag.RSSIWindow[0].PeakRssi
	maxIndex := 0

	for i, entry := range tag.RSSIWindow {
		if entry.PeakRssi >= maxVal {
			maxVal = entry.PeakRssi
			maxIndex = i
		}
	}

	// Ensure a significant drop after the peak
	if maxIndex < len(tag.RSSIWindow)-1 {
		if maxVal-tag.RSSIWindow[maxIndex+1].PeakRssi >= PeakDropThreshold {
			tag.Endtime = tag.RSSIWindow[maxIndex].FirstSeenTimestamp // Update end time
			return true
		}
	}

	return false
}

func (tag *Tag) handleData() {

	for {
		select {
		case data := <-tag.ChDataFromReader:
			if (len(tag.Data) == 0) && (data.AntennaPort == 17) {
				tag.AddPortFlag = true
				ChLEDToUI <- LEDServerToUI{tag.EPC, "GREEN"}
				tag.LED = "GREEN"
			}
			if data.AntennaPort == 17 && tag.AddPortFlag {
				tag.Data = append(tag.Data, data)
			}
			if data.AntennaPort == 9 && tag.AddPortFlag {
				tag.Data = append(tag.Data, data)
				if len(tag.countPortNumType()) > 1 {
					tag.maxIndexPower()
					tag.Dist = Distance
					fmt.Println("tag.Starttime", tag.Starttime)
					tag.Data = []RFIDData{}
					tag.RSSIWindow = []RSSIEntry{} // clear window for the second antenna
				}

				tag.updateWindow(data)
				if tag.detectLastPeak() {
					timeEnd := data.FirstSeenTimestamp
					fmt.Println("timeEnd", timeEnd)
					fmt.Println("tag.Endtime", tag.Endtime)
					timeDiff := float64(timeEnd-tag.Starttime) / 1000000
					speed := float64(Distance / timeDiff)
					speed = math.Round(speed*1000) / 1000
					currentTime := time.Now()
					timeText := currentTime.Format("15:04:05 PM")
					ChLEDToUI <- LEDServerToUI{tag.EPC, "RED"}
					tag.LED = "RED"
					ChDataToUI <- DataServerToUI{tag.EPC, speed, timeText, tag.Dist}
					fmt.Println("speed:", tag.EPC, speed, Distance, timeDiff)

					// reset all values
					tag.Data = []RFIDData{}
					tag.RSSIWindow = []RSSIEntry{}
					tag.Starttime = 0
					tag.Endtime = 0
					timeEnd = 0
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
