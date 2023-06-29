package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AmyangXYZ/sgo"
	"github.com/AmyangXYZ/sgo/middlewares"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	TagHolder          = make(map[string]*Tag)
	Distance   float64 = 0.86 //unit m
	ChDataToUI         = make(chan DataServerToUI, 1024)
	ChLEDToUI          = make(chan LEDServerToUI, 1024)
)

// sever to UI

type LEDServerToUI struct { //go AaaBbb //json: zz_zz //
	Epc string `json:"epc"`
	Led string `json:"led"`
}

type DataServerToUI struct { //go AaaBbb //json: zz_zz //
	Epc      string  `json:"epc"`
	CalSpeed float64 `json:"gait_speed"`
	Time     string  `json:"time"`
	Led      string  `json:led`
}

// reader to server
type RFIDData struct { //retrive data from RFID reader with format{epc  antennaport  firstseentime stamp}
	Epc                string `json:"epc"`
	AntennaPort        int    `json:"antennaPort"`
	FirstSeenTimestamp int64  `json:"firstSeenTimeStamp"`
}

type RFIDDataFromReader struct {
	TagReads []RFIDData `json:"tag_reads"` //?
}

func main() {

	app := sgo.New()
	app.USE(middlewares.CORS(middlewares.CORSOpt{}))

	app.GET("/", func(ctx *sgo.Context) error {
		return ctx.Text(200, "hello123")
	})

	app.POST("/api/reader/connect", PostFromReader) // receive data from rfid reader

	app.GET("/api/ui/tag", GetAllTags)                   // retrieve tag list (or tag holder)
	app.POST("/api/ui/tag/:id", PostTag)                 // register a tag by name
	app.OPTIONS("/api/ui/tag/:id", sgo.PreflightHandler) // handle CORS??
	app.GET("/api/ui/ws", GetWebSocket)                  // server data to UI
	app.Run(":16311")                                    //block
}

func PutDistance(ctx *sgo.Context) error {
	Distance, _ = strconv.ParseFloat(ctx.Param("distance"), 64)
	return ctx.Text(200, fmt.Sprintf("current distance is %v", Distance))
}

func PostFromReader(ctx *sgo.Context) error {
	// json body request
	body, err := io.ReadAll(ctx.Req.Body)
	if err != nil {
		return err
	}
	var data RFIDDataFromReader // raw read struct
	json.Unmarshal(body, &data)
	fmt.Println(data) //len could be 1 or 2 // use first one

	readerEpcInput24 := data.TagReads[0].Epc
	if _, ok := TagHolder[readerEpcInput24]; ok { //only pass data if key exist
		// if TagHolder[readerEpcInput24].AddPortFlag  // each tag will check this flag itself
		TagHolder[readerEpcInput24].ChDataFromReader <- data.TagReads[0]
	} else {
		fmt.Println("please register Tag")
	}
	return ctx.Text(200, "got it")
}

func GetAllTags(ctx *sgo.Context) error {
	return ctx.JSON(200, 1, "success", TagHolder)
}

func PostTag(ctx *sgo.Context) error { //from UI to server register new tag
	// id = "00000945" or "18145536"
	id := ctx.Param("id") // name must be in "epc{number}" format
	fmt.Println(id)
	id24 := strings.Repeat("0", 24-len(id)) + id
	epc := "epc" + id
	tag := newTag(epc, id24)
	ChLEDToUI <- LEDServerToUI{epc, "grey"}
	TagHolder[id24] = tag
	fmt.Println(id24)
	go tag.handleData()

	fmt.Println("TagHolder", TagHolder)
	return ctx.JSON(200, 1, "success", TagHolder)
}

func GetWebSocket(ctx *sgo.Context) error {
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
	count := 0
	for {
		select {
		// case data := <-ChDataToUI:
		// 	ws.WriteJSON(data)

		case <-time.After(1 * time.Second):
			count++
			if count%7 == 0 {

				hour, min, sec := time.Now().Clock()
				timeText := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)
				data := DataServerToUI{"18145536", 22.22, timeText, "RED"}
				ws.WriteJSON(data)

			} else if count%3 == 0 {
				// data := LEDServerToUI{"18145536", "GREEN"}
				tag := TagHolder["000000000000000000000945"]
				data := LEDServerToUI{tag.EPC, tag.LED}
				ws.WriteJSON(data)

			} else if count%4 == 0 {
				data := LEDServerToUI{"18145536", "RED"}
				ws.WriteJSON(data)
			}

		// case data := <-ChLEDToUI:
		// 	ws.WriteJSON(data)

		case <-breakSig:
			return nil
		}
	}
}
