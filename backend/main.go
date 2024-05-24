package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
	Distance   float64 = readDistanceFromFile("distconfig.txt") //unit m
	ChDataToUI         = make(chan DataServerToUI, 1024)
	ChLEDToUI          = make(chan LEDServerToUI, 1024)
	ChBreak            = make(chan bool, 1024)

	// TagList = []string{"E28068940000401D6E135DC6", "E28068940000501D6E13EDC6", "E28068940000401D6E13B5C6", "E28068940000401D6E136DC6", "000000000000000018145536"}
	TagList = map[string]string{
		"Tag1": "E28068940000401D6E135DC6",
		"Tag2": "E28068940000501D6E13EDC6",
		"Tag3": "E28068940000401D6E13B5C6",
		"Tag4": "E28068940000401D6E136DC6",
		"Tag5": "000000000000000000001088",
	}
)

// hard coded tags

// sever to UI
type LEDServerToUI struct {
	Epc string `json:"epc"`
	Led string `json:"led"`
}

type DataServerToUI struct {
	Epc      string  `json:"epc"`
	CalSpeed float64 `json:"gait_speed"`
	Time     string  `json:"time"`
	Dist     float64 `json:"dist"`
}

// reader to server
type RFIDData struct {
	Epc                string `json:"epc"`
	AntennaPort        int    `json:"antennaPort"`
	FirstSeenTimestamp int64  `json:"firstSeenTimeStamp"`
	PeakRssi           int    `json:"peakRssi"`
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
	app.POST("/api/reader/connect", PostFromReader) // recfeive data from rfid reader
	app.GET("/api/ui/tag", GetAllTags)              // retrieve tag list (or tag holder)

	app.POST("/api/ui/tag/:id", PostTag) // register a tag by name
	app.POST("/api/ui/tagdelete/:id24", DeleteTag)
	app.POST("/api/ui/distance/:UIdist", SetDistance)
	app.OPTIONS("/api/ui/tag/:id", sgo.PreflightHandler) // handle CORS
	app.GET("/api/ui/ws", GetWebSocket)                  // server data to UI
	AddHardCodeTag(TagList)
	app.Run(":16311") //block 16311

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
	// readerEpcInput24 := data.TagReads[0].Epc
	fmt.Println("---", len(data.TagReads)) //for debug
	for _, readerdatasingle := range data.TagReads {
		readerEpcInput24 := readerdatasingle.Epc
		fmt.Println("---", readerdatasingle)          //for debug
		if _, ok := TagHolder[readerEpcInput24]; ok { //only pass data if key exist
			// if TagHolder[readerEpcInput24].AddPortFlag  // each tag will check this flag itself
			TagHolder[readerEpcInput24].ChDataFromReader <- data.TagReads[0]
		} else {
			fmt.Println("unregistered tag is dectected")
		}
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

func SetDistance(ctx *sgo.Context) error {
	Distance, _ = strconv.ParseFloat(ctx.Param("UIdist"), 64)

	fmt.Println("set current distance", Distance)
	// return ctx.JSON(200, 1, "success", TagHolder)

	return ctx.Text(200, fmt.Sprintf("Distance set to %v meters", Distance))
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

func AddHardCodeTag(inList map[string]string) {
	for key, value := range inList {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
		id24 := strings.Repeat("0", 24-len(value)) + value
		// epc := "epc" + value
		epc := key
		tag := newTag(epc, id24)
		TagHolder[id24] = tag
		go tag.handleData()
	}

}
func readDistanceFromFile(filepath string) float64 {
	// Open the text file
	file, err := os.Open(filepath)
	if err != nil {
		return 0
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains "distance:"
		if strings.Contains(line, "distance:") {
			// Split the line by ":" to extract the value
			parts := strings.Split(line, ":")

			// Extract and convert the value to a float
			distanceStr := strings.TrimSpace(parts[1])
			distance, err := strconv.ParseFloat(distanceStr, 64)
			if err != nil {
				return 0
			}
			return distance
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return 0
	}

	return 0
}

func GetWebSocket(ctx *sgo.Context) error {
	fmt.Println("RFID ready to run ")

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
