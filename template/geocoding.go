package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

var (
	geometry Geocoding
	key      = os.Getenv("GEOCODING_API")
)

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

func main() {
	loc := "３１９御山村上門田町大字会津若松市福島県"
	url := endpoint + loc + "&key=" + key

	if err := GeometReq(url); err != nil {
		log.Print(err)
	}
	fmt.Println(geometry.Results[0].GeoRes.Location.Lat, geometry.Results[0].GeoRes.Location.Lng)
}
