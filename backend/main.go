package main

import (
	"fmt"
	"math/rand"
	"net/http"
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

func main() {
	app := sgo.New()
	app.GET("/", func(ctx *sgo.Context) error {
		return ctx.Text(200, "hello")
	})
	app.GET("/ws", ws)
	app.Run(":16311")
}

type RFIDRecord struct {
	Epc       int     `json:"epc"`
	Speed     float64 `json:"speed"`
	Timestamp int64   `json:"timestamp"`
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
	for {
		select {
		case <-time.After(1 * time.Second):
			record := RFIDRecord{epc, rand.Float64(), time.Now().Unix()}
			ws.WriteJSON(record)
			fmt.Println(record)
			epc++
		case <-breakSig:
			return nil
		}
	}
}
