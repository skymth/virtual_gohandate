package main

import (
	"math/rand"
	"time"
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

type buttonTemp struct {
	image string
	title string
	label string
}

type buttonTemp2 struct {
	key     string
	image   string
	title   string
	label   string
	select1 string
	select2 string
}

type buttonTemp4 struct {
	key     string
	title   string
	image   string
	label   string
	select1 string
	select2 string
	select3 string
	select4 string
}

var (
	button  map[string]buttonTemp
	button2 map[string]buttonTemp2
	button4 map[string]buttonTemp4
	talk    map[int]buttonTemp4
	word    map[string]string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	button = map[string]buttonTemp{
		"いただきます！": buttonTemp{
			image: "https://i.imgur.com/97XRjTa.png",
			title: "いただきます♪",
			label: " ",
		},
	}

	button2 = map[string]buttonTemp2{
		"ご飯行かない？": buttonTemp2{
			key:     "meshi",
			image:   "https://i.imgur.com/tNxL35o.png",
			title:   "しおえってお店行ってみたいなぁ♪",
			label:   " ",
			select1: "行く",
			select2: "行かない",
		},
		"location": buttonTemp2{
			key:     "osusume",
			image:   "https://i.imgur.com/AZ9L8d6.png",
			title:   "待ってたよ〜〜、お腹すいちゃった😂",
			label:   "オススメを聞く？聞かない？",
			select1: "聞く",
			select2: "聞かない",
		},
	}

	//	button4 = map[string]buttonTemp4{
	//		"": buttonTemp4{
	//			key:     "",
	//			title:   "",
	//			image:   "",
	//			label:   "",
	//			select1: "",
	//			select2: "",
	//			select3: "",
	//			select4: "",
	//		},
	//		"": buttonTemp4{
	//			key:     "",
	//			title:   "",
	//			image:   "",
	//			label:   "",
	//			select1: "",
	//			select2: "",
	//			select3: "",
	//			select4: "",
	//		},
	//	}

	talk = map[int]buttonTemp4{
		0: buttonTemp4{
			key:     "shiro",
			title:   "なんのお話にする？",
			image:   "https://i.imgur.com/iazlG5a.png",
			label:   " ",
			select1: "大阪城？",
			select2: "鶴ヶ城？",
			select3: "名古屋城？",
			select4: "カリオストロの城？",
		},
		1: buttonTemp4{
			key:     "kyodo",
			title:   "なんの話にする？",
			image:   "https://i.imgur.com/iazlG5a.png",
			label:   " ",
			select1: "ちゃんちゃんやき", //北海道の郷土料理だけどあんまわかんないや
			select2: "かにまき汁",    // 宮崎県の郷土料理だけどあんまわかんないや
			select3: "イノシシカレー",  // 山梨県の郷土料理だけどあんまわかんないや
			select4: "こづゆ",
		},
		2: buttonTemp4{
			key:     "men",
			title:   "なんの話にする？",
			image:   "https://i.imgur.com/iazlG5a.png",
			label:   " ",
			select1: "喜多方ラーメン",
			select2: "白河ラーメン",
			select3: "博多ラーメン",
			select4: "札幌ラーメン",
		},
	}

	word = map[string]string{
		"location": "嘘つき！\n全然違う場所じゃない！！",
		"meshi1":   "やったぁ！\n着いたら教えてね♪",
		"meshi2":   "そっかぁ...残念。。。",
		"osusume2": "了解♪",
		"menu1":    "おぉ！良いね♪\n私もソースカツ丼にしよう！",
		"menu2":    "おぉ！良いね♪\n私も味噌チャーシューにしよう！",
	}

}
