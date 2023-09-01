package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	ChBreak            = make(chan bool, 1024)
)

// sever to UI
type LEDServerToUI struct {
	Epc string `json:"epc"`
	Led string `json:"led"`
}

type DataServerToUI struct {
	Epc      string  `json:"epc"`
	CalSpeed float64 `json:"gait_speed"`
	Time     string  `json:"time"`
}

// reader to server
type RFIDData struct { //retrive data from RFID reader with format{epc  antennaport  firstseentime stamp}
	Epc                string `json:"epc"`
	AntennaPort        int    `json:"antennaPort"`
	FirstSeenTimestamp int64  `json:"firstSeenTimeStamp"`
}

type RFIDDataFromReader struct {
	TagReads []RFIDData `json:"tag_reads"`
}

func main() {
	app := sgo.New()
	app.SetTemplates("templates", nil)
	app.USE(middlewares.CORS(middlewares.CORSOpt{}))
	app.GET("/", func(ctx *sgo.Context) error {
		return ctx.Render(200, "i")
	})
	app.GET("/assets/*files", assets)
	app.POST("/api/reader/connect", PostFromReader) // receive data from rfid reader
	app.GET("/api/ui/tag", GetAllTags)              // retrieve tag list (or tag holder)
	app.POST("/api/ui/tag/:id", PostTag)            // register a tag by name
	app.POST("/api/ui/tagdelete/:id24", DeleteTag)
	app.OPTIONS("/api/ui/tag/:id", sgo.PreflightHandler) // handle CORS
	app.GET("/api/ui/ws", GetWebSocket)                  // server data to UI
	app.Run(":16311")                                    //block
}

// Static files handler.
func assets(ctx *sgo.Context) error {
	staticHandle := http.StripPrefix("/assets",
		http.FileServer(http.Dir("./assets")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
	return nil
}

func PostFromReader(ctx *sgo.Context) error {
	// json body request
	body, err := io.ReadAll(ctx.Req.Body)
	if err != nil {
		return err
	}
	var data RFIDDataFromReader
	json.Unmarshal(body, &data)
	readerEpcInput24 := data.TagReads[0].Epc
	fmt.Println(readerEpcInput24)
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

func DeleteTag(ctx *sgo.Context) error {
	id24 := ctx.Param("id24")

	if _, ok := TagHolder[id24]; ok {
		TagHolder[id24].ChSigBreak <- true
		delete(TagHolder, id24)
	}
	return ctx.JSON(200, 1, "success", nil)
}

func PostTag(ctx *sgo.Context) error { //from UI to server register new tag
	id := ctx.Param("id")
	id24 := strings.Repeat("0", 24-len(id)) + id
	epc := "epc" + id
	tag := newTag(epc, id24)
	TagHolder[id24] = tag
	go tag.handleData()
	return ctx.JSON(200, 1, "success", nil)
}

func GetWebSocket(ctx *sgo.Context) error {
	fmt.Println("Refresh GetWebSocket ")
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
	for {
		select {
		case data := <-ChDataToUI:
			ws.WriteJSON(data)
		case ledData := <-ChLEDToUI:
			ws.WriteJSON(ledData)
		case <-breakSig:
			return nil
		}
	}
}
