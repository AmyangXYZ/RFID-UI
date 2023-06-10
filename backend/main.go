package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/AmyangXYZ/sgo"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type RFIDDataRead struct {
	AntennaRead []dataRead `json:"tag_reads"` //?
}
type dataRead struct {
	Epc                string `json:"epc"`
	AntennaPort        int    `json:"antennaPort"`
	FirstSeenTimestamp int64  `json:"firstSeenTimeStamp"`
}
type TagInfo struct {
	taginfoEPC         string
	taginfoPort        []int
	taginfoTimeStamp   []int64
	taginfoAddPortFlag bool
}

var TagHolder = make(map[string]TagInfo)
var distance float64 = 0.86 //unit m //frontend set distance

func main() {

	// register button in frontend send to backend
	userInput := []string{"00000945", "18145536"}

	for _, userepc := range userInput {
		taginfoName := "epc" + userepc
		userInput24 := strings.Repeat("0", 24-len(userepc)) + userepc
		TagHolder[userInput24] = TagInfo{}
		//send "grey LED" to front end
		if mapvalue, ok := TagHolder[userInput24]; ok { // only run if  when ok == true
			mapvalue.taginfoAddPortFlag = true
			mapvalue.taginfoEPC = taginfoName //modify the copy
			TagHolder[userInput24] = mapvalue //reassign map entry
		}
	}

	app := sgo.New()
	app.GET("/", func(ctx *sgo.Context) error {
		return ctx.Text(200, "hello")
	})

	app.POST("/connect", rxRFIDData) //retrive data from rfid reader

	app.GET("/ws", ws)
	app.Run(":16311")
}

func rxRFIDData(ctx *sgo.Context) error {
	// json body request
	body, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		return err
	}
	var data RFIDDataRead // raw read struct
	json.Unmarshal(body, &data)
	// fmt.Println(data.AntennaRead) //len could be 1 or 2 // use first one

	postEPC := data.AntennaRead[0].Epc
	postPort := data.AntennaRead[0].AntennaPort
	postTime := data.AntennaRead[0].FirstSeenTimestamp

	if value, ok := TagHolder[postEPC]; ok { //if exist and AddPortFlag
		if value.taginfoAddPortFlag {
			//add port into port array
			value.taginfoPort = append(value.taginfoPort, postPort)
			value.taginfoTimeStamp = append(value.taginfoTimeStamp, postTime)
			TagHolder[postEPC] = value //reassign
		}
	} else {
		fmt.Println("tag not found, please register")
	}

	// iterate map TagHolder
	for key, value := range TagHolder {
		fmt.Println("iterate", key, value)
		if len(countPortNumType(TagHolder[key].taginfoPort)) == 1 {
			fmt.Println(key, len(countPortNumType(TagHolder[key].taginfoPort)))
		} else if len(countPortNumType(TagHolder[key].taginfoPort)) == 2 {
			ch := make(chan bool)

			first := TagHolder[key].taginfoTimeStamp[0]
			last := TagHolder[key].taginfoTimeStamp[len(TagHolder[key].taginfoTimeStamp)-1]
			timeDiff := float64(last-first) / 1000000
			speed := float64(distance / timeDiff)
			fmt.Println("speed:", key, speed)
			//send frontend "Red LED"

			value.taginfoAddPortFlag = false
			value.taginfoPort = make([]int, 0)
			value.taginfoTimeStamp = make([]int64, 0)
			TagHolder[key] = value

			go HoldTimer(ch, timeDiff, key)
			ch <- true
		}
	}

	//iterate map, if any register tag has countPortNumType len ==1
	//send frontend 'Green LED'

	//if countPortNumType == 2,
	//calcuate speed, send frontend epc:speed
	//set flag = false, send frontend 'Red LED'
	//set timer = 2 * time

	//time done
	//time up then flag = true
	//send fronfend "grey LED"

	// if len(data.AntennaRead) > 0 {
	// 	ch <- data
	// }

	return ctx.Text(200, "got it")
}

func ws(ctx *sgo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Resp, ctx.Req, nil)
	breakSig := make(chan bool)
	if err != nil {
		return err
	}

	defer func() {
		ws.Close()
	}()
	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				breakSig <- true
			}
		}
	}()
	epc := 0
	fmt.Println("dummy", epc)
	for {
		select {
		case <-time.After(1 * time.Second):
			//recond should be {epc speed} from rxRFIDData
			// record := dataRead{epc, rand.Float64(), time.Now().Unix()}
			// ws.WriteJSON(record)
			// fmt.Println(record)
			// epc++
		case <-breakSig:
			return nil
		}
	}
}

func countPortNumType(arr []int) map[int]int { //return number of type of ports
	var dict = make(map[int]int)
	for _, num := range arr {
		dict[num] = dict[num] + 1
	}

	return dict
}

func HoldTimer(ch chan bool, timeDuration float64, inKey string) {
	tmp := <-ch
	if tmp {
		fmt.Println("HoldTimer start ====== ", timeDuration)
		time.Sleep(time.Duration(2*timeDuration) * time.Second)
		fmt.Println("HoldTimer end ======", timeDuration*2)
		if value, ok := TagHolder[inKey]; ok {
			value.taginfoAddPortFlag = true
			TagHolder[inKey] = value
		}
	}
	close(ch)
	fmt.Println("channel close")
}
