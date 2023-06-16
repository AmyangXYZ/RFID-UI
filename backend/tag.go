package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Tag struct {
	ChDataFromReader chan RFIDData `json:"-"`
	ChSigBreak       chan bool     `json:"-"` // use - to omit it when convert to JSON
	ID24             string        `json:"id24"`
	EPC              string        `json:"epc"`
	Data             []RFIDData    `json:"data"`
	AddPortFlag      bool          `json:"add_port_flag"`
}

func newTag(id24 string) *Tag {
	return &Tag{
		ChDataFromReader: make(chan RFIDData, 1024),
		ChSigBreak:       make(chan bool),
		ID24:             id24,
		EPC:              "epc" + id24,
		AddPortFlag:      true,
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
				tag.Data = append(tag.Data, data)
				// tag.Ports = append(tag.Ports, postPort)
				// tag.Timestamps = append(tag.Timestamps, postTime)
				fmt.Println(tag.EPC, tag.Data)
			}
			if len(tag.countPortNumType()) > 1 {
				tag.AddPortFlag = false // no need a setter, since this attribute is public
				timeRangeStart := tag.Data[0].FirstSeenTimestamp
				timeRangeEnd := tag.Data[len(tag.Data)-1].FirstSeenTimestamp
				timeDiff := float64(timeRangeEnd-timeRangeStart) / 1000000
				speed := float64(Distance / timeDiff)
				speed = math.Round(speed*1000) / 1000
				hour, min, sec := time.Now().Clock()

				timeText := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)

				ChDataToUI <- DataServerToUI{tag.EPC, speed, timeText}
				fmt.Println("speed:", tag.EPC, speed)
				//send frontend "Red LED"

				tag.Data = []RFIDData{}

				fmt.Println(tag.EPC, "timer start")
				go func() {
					time.Sleep(time.Duration(2*timeDiff) * time.Second)
					tag.AddPortFlag = true
					fmt.Println(tag.EPC, "timer end, ready to run")
				}()

			}
		case <-tag.ChSigBreak:
			fmt.Println("quit", tag.EPC)
			return
		}
	}
}
