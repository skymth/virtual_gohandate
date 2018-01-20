package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	endpoint = "https://maps.googleapis.com/maps/api/geocode/json?address="
)

type User struct {
	date     bool
	location bool
}

type Geocoding struct {
	Results []Geometry `json:"results"`
}

type Geometry struct {
	GeoRes Location `json:"geometry"`
}

type Location struct {
	Location locations `json:"location"`
}

type locations struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Locs struct {
	Lat float64
	Lng float64
}

var (
	bot      *linebot.Client
	messages []linebot.Message
	geometry Geocoding
	key      = os.Getenv("GEOCODING_API")
	user     User
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
				case "ご飯行かない？":
					if _, err = bot.ReplyMessage(event.ReplyToken,
						linebot.NewTextMessage("行きたい！行きたい！"),
						ButtonTemplate2("https://i.imgur.com/tNxL35o.png", //<-変更
							"行く",
							"行かない",
							"しおえってお店行ってみたいなぁ♪",
							" ",
							"1",
							"2",
						)).Do(); err != nil {
						log.Print(err)
					}
				case "いただきます！":
					if _, err = bot.ReplyMessage(event.ReplyToken,
						ButtonTemplate("https://i.imgur.com/97XRjTa.png", //<-変更
							"いただきます♪",
						)).Do(); err != nil {
						log.Print(err)
					}

				case "ごちそうさま！":
					if _, err = bot.ReplyMessage(event.ReplyToken,
						linebot.NewTextMessage("ごちそうさまでした"),
						ReviewTemplate("https://i.imgur.com/oxoKeI5.png"), //<-変更
					).Do(); err != nil {
						log.Print(err)
					}

				case "お話ししよう！":
				default:
				}
			case *linebot.LocationMessage:
				LocationRes(message, event)
			}

		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "1":
				user.date = true
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("やったぁ！"),
					linebot.NewTextMessage("ついたら教えてね♪"),
				).Do(); err != nil {
					log.Print(err)
				}

			case "2":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("そっかぁ...残念。。。")).Do(); err != nil {
					log.Print(err)
				}

			case "3":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					CarouselTemplate()).Do(); err != nil {
					log.Print(err)
				}

			case "4":
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("了解♪")).Do(); err != nil {
					log.Print(err)
				}

			case "5":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("おぉ！良いね♪\n私もソースカツ丼にしよう！"),
				).Do(); err != nil {
					log.Print(err)
				}
			case "6":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("おぉ！良いね♪\n私も味噌チャーシューにしよう！"),
				).Do(); err != nil {
					log.Print(err)
				}
			case "7":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("また、来ようね！！"),
				).Do(); err != nil {
					log.Print(err)
				}
			default:
			}
		}
	}
}

func CarouselTemplate() *linebot.TemplateMessage {
	template := linebot.NewCarouselTemplate(
		linebot.NewCarouselColumn(
			"https://img.retty.me/img_repo/l/01/11709105.jpg",
			"味の里　しおえ",
			"ソースカツ丼",
			linebot.NewPostbackTemplateAction("これにする!", "5", ""),
		),
		linebot.NewCarouselColumn(
			"https://i.imgur.com/9Oam9dS.jpg",
			"味の里　しおえ",
			"味噌チャーシュー",
			linebot.NewPostbackTemplateAction("これにする!", "6", ""),
		),
	)

	msg := linebot.NewTemplateMessage("carousel", template)
	return msg
}

func ButtonTemplate(image, title string) *linebot.TemplateMessage {
	template := linebot.NewButtonsTemplate(
		image, //not image
		title, //ButtonsTemplate Title
		" ",   //ButtonsTemplate SubTitle
		linebot.NewPostbackTemplateAction(" ", " ", ""),
	)
	msg := linebot.NewTemplateMessage("confilm", template)
	return msg
}

func ButtonTemplate2(image, rb, lb, label, sublabel, no1, no2 string) *linebot.TemplateMessage {
	template := linebot.NewButtonsTemplate(
		image,    //not image
		label,    //ButtonsTemplate Title
		sublabel, //ButtonsTemplate SubTitle
		linebot.NewPostbackTemplateAction(rb, no1, ""),
		linebot.NewPostbackTemplateAction(lb, no2, ""),
	)
	msg := linebot.NewTemplateMessage("confilm", template)
	return msg
}

func ReviewTemplate(image string) *linebot.TemplateMessage {
	template := linebot.NewButtonsTemplate(
		image, //not image
		"美味しかったね！\n味はどうだった？", //ButtonsTemplate Title
		" ", //ButtonsTemplate SubTitle
		linebot.NewPostbackTemplateAction("★☆☆☆", "7", ""),
		linebot.NewPostbackTemplateAction("★★☆☆", "7", ""),
		linebot.NewPostbackTemplateAction("★★★☆", "7", ""),
		linebot.NewPostbackTemplateAction("★★★★", "7", ""),
	)
	msg := linebot.NewTemplateMessage("confilm", template)
	return msg
}

func GetConfirmData(rb, lb, label string) *linebot.TemplateMessage {
	rbutton := linebot.NewPostbackTemplateAction(rb, "1", "")
	lbutton := linebot.NewPostbackTemplateAction(lb, "2", "")

	temp := linebot.NewConfirmTemplate(label, rbutton, lbutton)
	msg := linebot.NewTemplateMessage("confilm-gohan", temp)

	return msg
}

func LocationRes(message *linebot.LocationMessage, event *linebot.Event) {
	loc := "３１９御山村上門田町大字会津若松市福島県"
	url := endpoint + loc + "&key=" + key
	if err := GeometReq(url); err != nil {
		log.Print(err)
	}

	la := strconv.FormatFloat(message.Latitude, 'f', 6, 64)
	lo := strconv.FormatFloat(message.Longitude, 'f', 6, 64)
	lat, _ := strconv.ParseFloat(la, 64)
	lon, _ := strconv.ParseFloat(lo, 64)

	max := Locs{}
	min := Locs{}
	max.Lat = geometry.Results[0].GeoRes.Location.Lat + 0.0040000
	max.Lng = geometry.Results[0].GeoRes.Location.Lng + 0.0020000

	min.Lat = geometry.Results[0].GeoRes.Location.Lat - 0.0040000
	min.Lng = geometry.Results[0].GeoRes.Location.Lng - 0.0020000

	if (lat >= min.Lat && lat <= max.Lat) && (lon >= min.Lng && lon <= max.Lng) {
		if _, err := bot.ReplyMessage(event.ReplyToken,
			linebot.NewTextMessage("待ってたよ〜〜"),
			ButtonTemplate2("https://i.imgur.com/AZ9L8d6.png", //<-変更
				"聞く",
				"聞かない",
				"待ってたよ〜〜、お腹すいちゃった😂",
				"オススメを聞く？聞かない？",
				"3",
				"4",
			),
		).Do(); err != nil {
			log.Print(err)
		}
	} else {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("嘘つき！\n全然違う場所じゃない！！")).Do(); err != nil {
			log.Print(err)
		}
	}

}

func GeometReq(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &geometry)
	if err != nil {
		return err
	}

	return nil
}
