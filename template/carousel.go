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
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							"https://img.retty.me/img_repo/l/01/11709105.jpg",
							"味の里　しおえ",
							"ソースカツ丼",
							linebot.NewPostbackTemplateAction("これにする!", "1", ""),
						),
						linebot.NewCarouselColumn(
							"https://img.retty.me/img_repo/l/01/1361436.jpg",
							"味の里　しおえ",
							"ロースカツ定食",
							linebot.NewPostbackTemplateAction("これにする!", "2", ""),
						),
					)

					msg := linebot.NewTemplateMessage("carousel", template)

					if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
						log.Print(err)
					}
				default:
				}
			}

		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "1":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("さすが、やっぱりこれだよね！")).Do(); err != nil {
					log.Print(err)
				}
			case "2":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("わかる〜〜！これこれ！")).Do(); err != nil {
					log.Print(err)
				}
			default:
			}
		}
	}
}
