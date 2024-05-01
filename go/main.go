package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

var (
	compass_brackets        = [17]string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW", "N"}
	OCEAN_COUNTRY    string = "Ocean"
	CITY_UNKNOWN     string = "Unknown"
)

const (
	earthRadiusMi = 3958 // radius of the earth in miles.
	earthRaidusKm = 6371 // radius of the earth in kilometers.
)

type issData struct {
	Message     string `json:"message"`
	Timestamp   int    `json:"timestamp"`
	IssPosition struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"iss_position"`
}

type groundLocation struct {
	PlaceID     int    `json:"place_id"`
	Licence     string `json:"licence"`
	PoweredBy   string `json:"powered_by"`
	OsmType     string `json:"osm_type"`
	OsmID       int    `json:"osm_id"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
	Address     struct {
		Village     string `json:"village"`
		County      string `json:"county"`
		State       string `json:"state"`
		Country     string `json:"country"`
		City        string `json:"city"`
		Suburb      string `json:"suburb"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
	Boundingbox []string `json:"boundingbox"`
}

type Coord struct {
	Lat float64
	Lon float64
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Distance calculates the shortest path between two coordinates on the surface
// of the Earth. This function returns two units of measure, the first is the
// distance in miles, the second is the distance in kilometers.
func Distance(p, q Coord) (mi, km float64) {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	mi = c * earthRadiusMi
	km = c * earthRaidusKm

	return mi, km
}

func direction_lookup(destination_x, destination_y, origin_x, origin_y float64) (compass string, deg int) {
	var degrees_final float64
	deltaX := destination_x - origin_x
	deltaY := destination_y - origin_y
	degrees_temp := math.Atan2(deltaX, deltaY) / math.Pi * 180
	if degrees_temp < 0 {
		degrees_final = 360 + degrees_temp
	} else {
		degrees_final = degrees_temp
	}
	compass_lookup := math.Round(degrees_final / 22.5)
	//fmt.Printf(strconv.Itoa(int(compass_lookup)))

	return compass_brackets[int(compass_lookup)], int(degrees_final)
}

func main() {

	DEVICE_LAT := 53.480970
	DEVICE_LON := -2.237150

	// Sign up for an account at geocode.maps.co for a free api key
	GEOCODE_API_KEY := ""

	country := OCEAN_COUNTRY
	city := CITY_UNKNOWN

	iss_raw, err := http.Get("http://api.open-notify.org/iss-now.json")

	if err != nil {
		log.Fatalln(err)
	}

	defer iss_raw.Body.Close()

	decoder := json.NewDecoder(iss_raw.Body)
	var issRaw issData
	err = decoder.Decode(&issRaw)
	if err != nil {
		log.Fatalln(err)
	}

	iss_lat, err := strconv.ParseFloat(issRaw.IssPosition.Latitude, 64)
	if err != nil {
		log.Fatalln(err)
	}

	iss_lon, err := strconv.ParseFloat(issRaw.IssPosition.Longitude, 64)
	if err != nil {
		log.Fatalln(err)
	}

	issLocation := Coord{Lat: iss_lat, Lon: iss_lon}
	DEVICE := Coord{Lat: DEVICE_LAT, Lon: DEVICE_LON}
	mi, km := Distance(issLocation, DEVICE)

	_ = km

	ground_url := "https://geocode.maps.co/reverse?lat=" + strconv.Itoa(int(iss_lat)) + "&lon=" + strconv.Itoa(int(iss_lon)) + "&api_key=" + string(GEOCODE_API_KEY)

	ground_raw, err := http.Get(ground_url)

	if err != nil {
		log.Fatalln(err)
	}

	defer ground_raw.Body.Close()

	grounddecoder := json.NewDecoder(ground_raw.Body)
	var groundRaw groundLocation
	err = grounddecoder.Decode(&groundRaw)
	if err != nil {
		log.Fatalln(err)
	}

	if groundRaw.Address.Country != "" {
		country = groundRaw.Address.Country
	}

	if country != OCEAN_COUNTRY {
		if groundRaw.Address.City != "" {
			city = groundRaw.Address.City
		}
		if groundRaw.Address.Suburb != "" {
			city = groundRaw.Address.Suburb
		}
		if groundRaw.Address.State != "" {
			city = groundRaw.Address.State
		}

	}

	fmt.Println("Iss Distance from your location is:", mi, "miles")
	fmt.Println("Country: ", country)
	fmt.Println("Nearest City:", city)

	direction, deg := direction_lookup(iss_lon, iss_lat, DEVICE_LON, DEVICE_LAT)
	fmt.Println("Viewing Direction:", direction, "(Bearing:", deg, ")")

}
