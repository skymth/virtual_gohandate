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
	geometry Geocoding
	key      = os.Getenv("GEOCODING_API")
	bot      *linebot.Client
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
	//geocoding
	loc := "３１９御山村上門田町大字会津若松市福島県"
	url := endpoint + loc + "&key=" + key
	if err := GeometReq(url); err != nil {
		log.Print(err)
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
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("近くにいるね")).Do(); err != nil {
						log.Print(err)
					}
				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("嘘つき！")).Do(); err != nil {
						log.Print(err)
					}
				}
			}

		default:
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("位置情報をお願いします。")).Do(); err != nil {
				log.Print(err)
			}
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
