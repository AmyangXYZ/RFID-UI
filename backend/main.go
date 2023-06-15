package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

// UI to server
type DataServertoUI struct { //go AaaBbb //json: zz_zz //
	Epc      string  `json:"epc"`
	CalSpeed float64 `json:"gait_speed"`
}

// reader to server
type DataTakeFromReader struct { //retrive data from RFID reader with format{epc  antennaport  firstseentime stamp}
	Epc                string `json:"epc"`
	AntennaPort        int    `json:"antennaPort"`
	FirstSeenTimestamp int64  `json:"firstSeenTimeStamp"`
}

type RFIDDataFromReader struct {
	TagReads []DataTakeFromReader `json:"tag_reads"` //?
}

// server to UI
var (
	TagHolder          = make(map[string]*Tag)
	distance   float64 = 0.86 //unit m
	chDataToUI         = make(chan DataServertoUI, 1024)
)

func newTag() *Tag { //? content in newtag??
	return &Tag{
		ChDataFromReader: make(chan RFIDDataFromReader, 1024),
		ChSigBreak:       make(chan bool),
		// sigReset
		// sigPause
		//
	}
}
func main() {

	// register button in frontend send to backend
	UIButtonTagInputArray := []string{"00000945", "18145536"}

	for _, UIButtonTagInput := range UIButtonTagInputArray {
		tagName := "epc" + UIButtonTagInput
		UIButtonTagInput24 := strings.Repeat("0", 24-len(UIButtonTagInput)) + UIButtonTagInput

		tag := newTag()
		// TagHolder[tag.id] = tag
		// //send "grey LED" to frontend
		TagHolder[UIButtonTagInput24] = tag
		TagHolder[UIButtonTagInput24].setEPCName(tagName)
		TagHolder[UIButtonTagInput24].setAddFlag(true)
		go tag.handleData()
	}
	fmt.Println("TagHolder", TagHolder)
	app := sgo.New()
	app.GET("/", func(ctx *sgo.Context) error {
		return ctx.Text(200, "hello")
	})

	app.POST("/api/reader/connect", rxRFIDData) //retrive data from rfid reader

	app.GET("/ws", ws)

	// // /api/ui/1/true
	// app.POST("/api/ui/:tag_id/:start_pause", pause)
	// // app.GET("/ws", ws)

	// app.POST("/api/ui/:tag_id", pause)              //
	// app.DELETE("/api/ui/:tag_id", deleteTagHandler) //

	// app.PUT("/api/ui/distance", updateDistance)
	app.Run(":16311") //block
}

// func deleteTagHandler(ctx *sgo.Context) {
// 	tag_id := ctx.Param("tag_id")                             //tag_id from UI is number only on physical tag
// 	tag_id_24 := strings.Repeat("0", 24-len(tag_id)) + tag_id //the key of map is 24 digit

// 	TagHolder[tag_id_24].ChSigBreak <- true
// 	delete(TagHolder, tag_id_24) // remove from map tagHolder[tag_id_24]

// }
// func updateDistance(ctx *sgo.Context) error {
// 	distance = ctx.Param("distance")
// 	return err
// }

func rxRFIDData(ctx *sgo.Context) error {
	// json body request
	body, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		return err
	}
	var data RFIDDataFromReader // raw read struct
	json.Unmarshal(body, &data)
	// fmt.Println(data) //len could be 1 or 2 // use first one

	readerEpcInput24 := data.TagReads[0].Epc
	if _, ok := TagHolder[readerEpcInput24]; ok { //only pass data if key exist
		if TagHolder[readerEpcInput24].TagAddPortFlag {
			TagHolder[readerEpcInput24].ChDataFromReader <- data
		}
	} else {
		fmt.Println("please register Tag")
	}
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

	// epc := 0
	// speed := 0.01
	for {
		select {
		// case <-time.After(1 * time.Second):
		// 	//recond should be {epc speed} from rxRFIDData
		// 	// record := DataTakeFromReader{"epc", epc, time.Now().Unix()}
		// 	record := DataServertoUI{"epc", speed}

		// 	ws.WriteJSON(record)
		// 	fmt.Println(record)
		// 	epc++
		// 	speed++
		case data := <-chDataToUI:
			ws.WriteJSON(data)
			fmt.Println("in ws function")

		case <-breakSig:
			return nil
		}
	}
}
