package main

import (
	"log"
	"net/http"
	"os"

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
			case *linebot.TextMessage:
				switch message.Text {
				case "test":
					template := linebot.NewButtonsTemplate(
						"",                 //not image
						"My button sample", //ButtonsTemplate Title
						"Hello, my button", //ButtonsTemplate SubTitle
						linebot.NewPostbackTemplateAction("banana", "1", ""),
						linebot.NewPostbackTemplateAction("設楽", "2", ""),
						linebot.NewPostbackTemplateAction("日村", "3", ""),
					)
					msg := linebot.NewTemplateMessage("confilm", template)

					if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
						log.Print(err)
					}
				default:
				}
			}

		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "1":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("man")).Do(); err != nil {
					log.Print(err)
				}
			case "2":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("統")).Do(); err != nil {
					log.Print(err)
				}
			case "3":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("勇気")).Do(); err != nil {
					log.Print(err)
				}
			default:
			}
		}
	}
}
