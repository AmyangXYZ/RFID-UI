package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Tag struct {
	ChDataFromReader chan RFIDDataFromReader
	ChSigBreak       chan bool
	InfoEPC          string  `json:"tag_info_epc"`
	TagPort          []int   `json:"tag_info_tagPort"`
	TagTimeStamp     []int64 `json:"tag_info_tagtimestamp"`
	TagAddPortFlag   bool    `json:"tag_info_tagaddportflag"`
}

func (tag *Tag) setEPCName(inEPC string) {
	tag.InfoEPC = inEPC
}

func (tag *Tag) setAddFlag(inBool bool) {
	tag.TagAddPortFlag = inBool
}

func (tag *Tag) countPortNumType(arr []int) map[int]int { //return {1:11, 9:1} ==> len=2
	var dict = make(map[int]int)
	for _, num := range arr {
		dict[num] = dict[num] + 1
	}

	return dict
}

func (tag *Tag) handleData() {
	for {
		select {

		case data := <-tag.ChDataFromReader:
			postPort := data.TagReads[0].AntennaPort
			postTime := data.TagReads[0].FirstSeenTimestamp
			if tag.TagAddPortFlag {
				tag.TagPort = append(tag.TagPort, postPort)
				tag.TagTimeStamp = append(tag.TagTimeStamp, postTime)
				fmt.Println(tag.InfoEPC, tag.TagPort)

			}
			if len(tag.countPortNumType(tag.TagPort)) == 2 {
				tag.setAddFlag(false)
				first := tag.TagTimeStamp[0]
				last := tag.TagTimeStamp[len(tag.TagTimeStamp)-1]
				timeDiff := float64(last-first) / 1000000
				speed := float64(distance / timeDiff)
				speed = math.Round(speed*1000) / 1000
				hour, min, sec := time.Now().Clock()

				timeText := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)

				chDataToUI <- DataServertoUI{tag.InfoEPC, speed, timeText}
				fmt.Println("speed:", tag.InfoEPC, speed)
				//send frontend "Red LED"

				tag.TagPort = make([]int, 0)
				tag.TagTimeStamp = make([]int64, 0)
				fmt.Println(tag.InfoEPC, "timer start")
				go func() {
					time.Sleep(time.Duration(2*timeDiff) * time.Second)
					tag.setAddFlag(true)
					fmt.Println(tag.InfoEPC, "timer end, ready to run")
				}()

			}
			break
		case <-tag.ChSigBreak:
			fmt.Println("quit", tag.InfoEPC)
			return
			break
		}
	}
	// for {
	// 	select {
	// 	case data := <-tag.ch:
	// 		if !tag.pause {

	// 			// Tag.TagPort = append(Tag.TagPort)

	// 			// if true {
	// 			// 	tag.pause = true
	// 			// 	go func() {
	// 			// 		time.sleep()
	// 			// 		tag.resumeSig <- true
	// 			// 	}()
	// 			// }
	// 		}
	// 		// if () {

	// 		// }
	// 		// 	break
	// 		// case <-tag.resumeSig:
	// 		// 	tag.pause = false
	// 		// 	break

	// 		// case <-tag.deleteSig:
	// 		// 	return
	// 		// 	break
	// 	}
	// }
}
