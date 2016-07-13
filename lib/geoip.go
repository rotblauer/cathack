package lib

import (
	// "fmt"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/oschwald/geoip2-golang"
)

func GetGeoFromIP(ip string) (out map[string]string, err error) {

	// fmt.Printf("Getting geo from ip: %v", ip)

	out = make(map[string]string)
	// Don't init error. Want it to be nil.

	db, err := geoip2.Open("./data/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// If you are using strings that may be invalid, check that ip is not nil
	if ip == "" {
		err = errors.New("IP is blank.")
		log.Fatal(err)
	}

	parsedIP := net.ParseIP(ip)

	record, err := db.City(parsedIP)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	// fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	// fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	// fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	// fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	// fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
	// Output:
	// Portuguese (BR) city name: Londres
	// English subdivision name: England
	// Russian country name: Великобритания
	// ISO country code: GB
	// Time zone: Europe/London
	// Coordinates: 51.5142, -0.0931
	out["lat"] = strconv.FormatFloat(record.Location.Latitude, 'f', 6, 64)
	out["lon"] = strconv.FormatFloat(record.Location.Longitude, 'f', 6, 64)
	out["tz"] = record.Location.TimeZone
	out["countryIsoCode"] = record.Country.IsoCode
	out["subdiv"] = record.Subdivisions[0].Names["en"]
	out["city"] = record.City.Names["en"]
	fmt.Printf("GeoIP: %v\n", out)
	fmt.Println()
	if err != nil {
		return out, err
	}
	return out, nil
}
