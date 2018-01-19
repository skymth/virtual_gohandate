package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot *linebot.Client
)

func main() {
	var err error
	bot, err = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", ResponseCall)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func ResponseCall(w http.ResponseWriter, req *http.Request) {
	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.LocationMessage:
				lat := strconv.FormatFloat(message.Latitude, 'f', 6, 64)
				lon := strconv.FormatFloat(message.Longitude, 'f', 6, 64)

				msg := fmt.Sprintf("緯度:%v\n経度:%v", lat, lon)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
					log.Print(err)
				}
			}

		default:
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("位置情報をお願いします。")).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}
