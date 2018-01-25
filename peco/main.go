package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	endpoint = "https://maps.googleapis.com/maps/api/geocode/json?address="
)

var (
	key = os.Getenv("GEOCODING_API")
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
			if err := peco.postbackResponse(event.Postback, event.ReplyToken); err != nil {
				log.Print(err)
			}
		default:
			log.Printf("Unknown: %v", event)
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
			return err
		}
	case "いただきます！":
		if _, err = peco.bot.ReplyMessage(
			reply,
			ButtonTemplate(button[message.Text]),
		).Do(); err != nil {
			return err
		}
	case "ごちそうさま！":
		if _, err = peco.bot.ReplyMessage(
			reply,
			linebot.NewTextMessage("ごちそうさまでした"),
			ButtonTemplate4(button4[message.Text], false),
		).Do(); err != nil {
			return err
		}
	case "お話しよう！":

		if _, err = peco.bot.ReplyMessage(
			reply,
			ButtonTemplate4(talk[rand.Intn(3)], true),
		).Do(); err != nil {
			return err
		}
	case "（好感度）":
		if err := peco.messageResponse(reply, formatStr("好感度: ", peco.favoRate)); err != nil {
			return err
		}
	default:
	}
	return nil
}

func (peco *Peco) locationResponse(message *linebot.LocationMessage, reply string) error {
	resLocation, typ, err := handleLocation(message.Latitude, message.Longitude)
	if err != nil {
		return err
	} else if typ == true {
		if _, err := peco.bot.ReplyMessage(
			reply,
			linebot.NewTextMessage(word["location"]),
		).Do(); err != nil {
			return err
		}
	} else {
		if _, err := peco.bot.ReplyMessage(
			reply,
			resLocation,
		).Do(); err != nil {
			return err
		}
	}
	return nil
}

func (peco *Peco) postbackResponse(p *linebot.Postback, reply string) error {
	switch p.Data {
	case "meshi1", "meshi2", "osusume2", "menu1", "menu2", "review":
		if err := peco.messageResponse(reply, word[p.Data]); err != nil {
			return err
		}
	case "osusume1":
		if err := peco.carouselResponse(reply); err != nil {
			return err
		}
	case "shiro1", "shiro2", "shiro3", "shiro4", "kyodo1", "kyodo2", "kyodo3", "kyodo4", "men1", "men2", "men3", "men4":
		if err := peco.talkResponse(reply, talkres[p.Data]); err != nil {
			return err
		}
	}
	return nil
}

func ButtonTemplate(res buttonTemp) *linebot.TemplateMessage {
	temp := linebot.NewButtonsTemplate(
		res.image,
		res.title,
		res.label,
		linebot.NewPostbackTemplateAction(res.select1, res.key, ""),
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

func ButtonTemplate4(res buttonTemp4, caset bool) *linebot.TemplateMessage {
	var temp *linebot.ButtonsTemplate
	if caset {
		temp = linebot.NewButtonsTemplate(
			res.image,
			res.title,
			res.label,
			linebot.NewPostbackTemplateAction(res.select1, formatStr(res.key, 1), ""),
			linebot.NewPostbackTemplateAction(res.select2, formatStr(res.key, 2), ""),
			linebot.NewPostbackTemplateAction(res.select3, formatStr(res.key, 3), ""),
			linebot.NewPostbackTemplateAction(res.select4, formatStr(res.key, 4), ""),
		)
	} else {
		temp = linebot.NewButtonsTemplate(
			res.image,
			res.title,
			res.label,
			linebot.NewPostbackTemplateAction(res.select1, res.key, ""),
			linebot.NewPostbackTemplateAction(res.select2, res.key, ""),
			linebot.NewPostbackTemplateAction(res.select3, res.key, ""),
			linebot.NewPostbackTemplateAction(res.select4, res.key, ""),
		)
	}
	return linebot.NewTemplateMessage("button4", temp)
}

func (peco *Peco) carouselResponse(reply string) error {
	var columns []*linebot.CarouselColumn
	columns = append(columns, NewColum(button["so-su"]))
	columns = append(columns, NewColum(button["miso"]))
	if _, err := peco.bot.ReplyMessage(
		reply,
		CarouselTemplate(columns),
	).Do(); err != nil {
		return err
	}
	return nil
}

func CarouselTemplate(columns []*linebot.CarouselColumn) *linebot.TemplateMessage {
	return linebot.NewTemplateMessage("carousel", linebot.NewCarouselTemplate(columns...))
}

func NewColum(res buttonTemp) *linebot.CarouselColumn {
	return linebot.NewCarouselColumn(
		res.image,
		res.title,
		res.label,
		linebot.NewPostbackTemplateAction(res.select1, res.key, ""),
	)
}

func handleLocation(lat, lon float64) (*linebot.TemplateMessage, bool, error) {
	loc := "３１９御山村上門田町大字会津若松市福島県" //example adress
	url := endpoint + loc + "&key=" + key
	geo, err := GeometReq(url)
	if err != nil {
		return nil, false, err
	}

	max := Locs{
		Lat: geo.Results[0].GeoRes.Location.Lat + 0.0040000,
		Lng: geo.Results[0].GeoRes.Location.Lng + 0.0020000,
	}
	min := Locs{
		Lat: geo.Results[0].GeoRes.Location.Lat - 0.0040000,
		Lng: geo.Results[0].GeoRes.Location.Lng - 0.0020000,
	}

	if (lat >= min.Lat && lat <= max.Lat) && (lon >= min.Lng && lon <= max.Lng) {
		return ButtonTemplate2(button2["location"]), false, nil
	} else {
		return nil, true, nil
	}

}

func formatStr(str string, i int) string {
	return fmt.Sprintf("%s%d", str, i)
}

func GeometReq(url string) (*Geocoding, error) {
	var geo Geocoding
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &geo)
	if err != nil {
		return nil, err
	}

	return &geo, nil
}

func (peco *Peco) messageResponse(replyToken, text string) error {
	if _, err := peco.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}

func (peco *Peco) talkResponse(reply string, res TalkRes) error {
	if res.image != "" {
		if _, err := peco.bot.ReplyMessage(
			reply,
			linebot.NewTextMessage(res.text),
			linebot.NewImageMessage(res.image, res.image),
		).Do(); err != nil {
			return err
		}
	} else {
		if _, err := peco.bot.ReplyMessage(
			reply,
			linebot.NewTextMessage(res.text),
			stickerRandom(),
		).Do(); err != nil {
			return err
		}
	}
	return nil
}

func stickerRandom() *linebot.StickerMessage {
	return linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1))
}
