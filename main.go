package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

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

type Template struct {
	Title   string
	image   string
	select1 string
	select2 string
	select3 string
	select4 string
	key     string
}

var (
	bot      *linebot.Client
	messages []linebot.Message
	geometry Geocoding
	key      = os.Getenv("GEOCODING_API")
	user     User
	comm     = map[int]Template{
		0: Template{
			Title:   "なんのお話にする？",
			image:   "https://i.imgur.com/iazlG5a.png",
			select1: "大阪城？",
			select2: "鶴ヶ城？",
			select3: "名古屋城？",
			select4: "カリオストロの城？",
			key:     "shiro",
		},
		1: Template{
			Title:   "なんの話にする？",
			image:   "https://i.imgur.com/iazlG5a.png",
			select1: "ちゃんちゃんやき", //北海道の郷土料理だけどあんまわかんないや
			select2: "かにまき汁",    // 宮崎県の郷土料理だけどあんまわかんないや
			select3: "イノシシカレー",  // 山梨県の郷土料理だけどあんまわかんないや
			select4: "こづゆ",
			key:     "kyodo",
		},
		2: Template{
			Title:   "なんの話にする？",
			image:   "https://i.imgur.com/iazlG5a.png",
			select1: "喜多方ラーメン",
			select2: "白河ラーメン",
			select3: "博多ラーメン",
			select4: "札幌ラーメン",
			key:     "men",
		},
	}
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

				case "お話しよう！":
					rand.Seed(time.Now().UnixNano())

					if _, err = bot.ReplyMessage(event.ReplyToken,
						SelectTemplate(comm[rand.Intn(3)]), //<-変更
					).Do(); err != nil {
						log.Print(err)
					}
				default:
				}
			case *linebot.LocationMessage:
				LocationRes(message, event)
			}

		case linebot.EventTypePostback:
			rand.Seed(time.Now().UnixNano())
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

			case "shiro1":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("それはあんま興味ないなぁ〜〜..."),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "shiro2":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("peco 鶴ヶ城には詳しいんだぁ〜！"),
					linebot.NewTextMessage("福島県会津若松市追手町にあった日本の城で、地元では鶴ヶ城（つるがじょう）と言うが、同名の城が他にあるため、地元以外では会津若松城と呼ばれることが多い。文献では旧称である黒川城（くろかわじょう）、または単に会津城とされることもある。国の史跡としては、若松城跡（わかまつじょうあと）の名称で指定されている。"),
					linebot.NewImageMessage("https://i.imgur.com/nPejtHV.jpg", "https://i.imgur.com/nPejtHV.jpg"),
				).Do(); err != nil {
					log.Print(err)
				}

			case "shiro3":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("名古屋城かぁ〜名古屋城はあんまり詳しくないんだぁ〜"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "shiro4":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("ルパ〜ン3世...だね！！！"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "kyodo1":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("北海道の郷土料理だけどあんまわかんないや"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "kyodo2":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("宮崎県の郷土料理だけどあんまわかんないや"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "kyodo3":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("山梨県の郷土料理だけどあんまわかんないや"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "kyodo4":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("内陸の会津地方でも入手が可能な、海産物の乾物を素材とした汁物である。江戸時代後期から明治初期にかけて会津藩の武家料理や庶民のごちそうとして広まり、現在でも正月や冠婚葬祭などハレの席で、必ず振る舞われる郷土料理である。"),
					linebot.NewImageMessage("https://i.imgur.com/uUWeU5G.jpg", "https://i.imgur.com/uUWeU5G.jpg"),
				).Do(); err != nil {
					log.Print(err)
				}

			case "men1":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("peco 喜多方ラーメン大好きなんだぁ！"),
					linebot.NewTextMessage("喜多方ラーメン（きたかたラーメン）とは福島県喜多方市発祥のご当地ラーメン（ご当地グルメ）で、2006年（平成18年）1月の市町村合併前の旧喜多方市では人口37,000人あまりに対し120軒ほどのラーメン店があり、対人口比の店舗数では日本一であった。札幌ラーメン、博多ラーメンと並んで日本三大ラーメンの一つに数えられている。"),
					linebot.NewImageMessage("https://i.imgur.com/w6kws4W.png", "https://i.imgur.com/w6kws4W.png"),
				).Do(); err != nil {
					log.Print(err)
				}

			case "men2":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("喜多方ラーメンの話ししようよ〜！！"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "men3":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("喜多方ラーメンの話ししようよ〜！！"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			case "men4":
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage("喜多方ラーメンの話ししようよ〜！！"),
					linebot.NewStickerMessage("1", fmt.Sprintf("%d", rand.Intn(18)+1)),
				).Do(); err != nil {
					log.Print(err)
				}

			default:
			}
		}
	}
}

func SelectTemplate(res Template) *linebot.TemplateMessage {
	template := linebot.NewButtonsTemplate(
		res.image, //not image
		res.Title, //ButtonsTemplate Title
		" ",       //ButtonsTemplate SubTitle
		linebot.NewPostbackTemplateAction(res.select1, fmt.Sprintf("%s1", res.key), ""),
		linebot.NewPostbackTemplateAction(res.select2, fmt.Sprintf("%s2", res.key), ""),
		linebot.NewPostbackTemplateAction(res.select3, fmt.Sprintf("%s3", res.key), ""),
		linebot.NewPostbackTemplateAction(res.select4, fmt.Sprintf("%s4", res.key), ""),
	)
	msg := linebot.NewTemplateMessage("confilm", template)
	return msg
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
