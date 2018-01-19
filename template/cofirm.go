package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot      *linebot.Client
	messages []linebot.Message
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
					//					rb := &linebot.MessageTemplateAction{
					//						Label: "right",
					//						Text:  "rrrrr",
					//					}
					//					lb := &linebot.MessageTemplateAction{
					//						Label: "left",
					//						Text:  "lllll",
					//					}
					rb := &linebot.PostbackTemplateAction{
						Label: "right",
						Data:  "rrrrr",
					}
					lb := &linebot.PostbackTemplateAction{
						Label: "left",
						Data:  "lllll",
					}

					temp := linebot.NewConfirmTemplate("aiueo", lb, rb)
					msg := linebot.NewTemplateMessage("confilm", temp)

					if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
						log.Print(err)
					}
				default:
				}
			}

		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "rrrrr":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("右が押されました")).Do(); err != nil {
					log.Print(err)
				}
			case "lllll":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("左が押されました")).Do(); err != nil {
					log.Print(err)
				}
			default:
			}
		}
	}
}
