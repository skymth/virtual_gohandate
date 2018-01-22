package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Peco struct {
	bot      *linebot.Client
	favoRate int
}

func main() {
	peco, err := NewClient()
	if err != nil {
		log.Print(err)
	}

	http.HandleFunc("/callback", peco.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func NewClient() (*Peco, error) {
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		return nil, err
	}

	return &Peco{
		bot: bot,
	}, nil
}

func (peco *Peco) Callback(w http.ResponseWriter, req *http.Request) {
	events, err := peco.bot.ParseRequest(req)
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
				if err := peco.textResponse(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				if err := peco.locationResponse(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			default:
			}

		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "":
			default:
			}
		}
	}
}

func (peco *Peco) textResponse(message *linebot.TextMessage, reply string) error {
	var err error
	switch message.Text {
	case "ご飯行かない？":
		if _, err = peco.bot.ReplyMessage(
			reply,
			linebot.NewTextMessage("行きたい！行きたい！"),
			ButtonTemplate2(button2[message.Text]),
		).Do(); err != nil {
			log.Print(err)
		}
	case "いただきます！":
		if _, err = peco.bot.ReplyMessage(
			reply,
			ButtonTemplate(button[message.Text]),
		).Do(); err != nil {
			log.Print(err)
		}
	case "ごちそうさま！":
		if _, err = peco.bot.ReplyMessage(
			reply,
			linebot.NewTextMessage("ごちそうさまでした"),
			ButtonTemplate4(button4[message.Text]),
		).Do(); err != nil {
			log.Print(err)
		}
	case "お話しよう！":
		if _, err = peco.bot.ReplyMessage(
			reply,
			ButtonTemplate4(talk[rand.Intn(3)]),
		).Do(); err != nil {
			log.Print(err)
		}
	case "（好感度）":
		if _, err = peco.bot.ReplyMessage(
			reply,
			linebot.NewTextMessage(formatStr("好感度: ", peco.favoRate)),
		).Do(); err != nil {
			log.Print(err)
		}
	default:
	}
	return nil
}

func (peco *Peco) locationResponse(message *linebot.LocationMessage, reply string) error {
	resLocation, err := handleLocation(message.Latitude, message.Longitude)
	if err != nil {
		return err
	}
	if _, err := app.bot.ReplyMessage(
		reply,
		resLocation,
	).Do(); err != nil {
		return err
	}
	return nil
}

func ButtonTemplate(res buttonTemp) *linebot.TemplateMessage {
	temp := linebot.NewButtonsTemplate(
		res.image,
		res.title,
		res.label,
		linebot.NewPostbackTemplateAction(" ", " ", ""),
	)
	return linebot.NewTemplateMessage("button", temp)
}

func ButtonTemplate2(res buttonTemp2) *linebot.TemplateMessage {
	temp := linebot.NewButtonsTemplate(
		res.image,
		res.title,
		res.label,
		linebot.NewPostbackTemplateAction(res.select1, formatStr(res.key, 1), ""),
		linebot.NewPostbackTemplateAction(res.select2, formatStr(res.key, 2), ""),
	)
	return linebot.NewTemplateMessage("button2", temp)
}

func ButtonTemplate4(res buttonTemp4) *linebot.TemplateMessage {
	temp := linebot.NewButtonsTemplate(
		res.image,
		res.title,
		res.label,
		linebot.NewPostbackTemplateAction(res.select1, formatStr(res.key, 1), ""),
		linebot.NewPostbackTemplateAction(res.select2, formatStr(res.key, 2), ""),
		linebot.NewPostbackTemplateAction(res.select3, formatStr(res.key, 3), ""),
		linebot.NewPostbackTemplateAction(res.select4, formatStr(res.key, 4), ""),
	)
	return linebot.NewTemplateMessage("confilm", temp)
}

func formatStr(str string, i int) string {
	return fmt.Sprintf("%s%d", str, i)
}
