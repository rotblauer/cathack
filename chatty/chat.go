package chatty

import (
	"os"
	"strconv"
	"time"
	"log"
	"fmt"

	"github.com/olahol/melody"
	j "github.com/ricardolonga/jsongo"

	"../lib"
)

func saveChat(data []byte) (bytes int, err error) {

	// Open database.
	f, err := os.OpenFile("./chat.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalln("Error opening file: ", err)
	}

	// Write to database.
	line := string(data)
	bytes, err = f.WriteString(line + "\n")
	if err != nil {
		log.Fatalln("Error writing string: ", err) // Will this out to same place as fmt? ie &>chat.log
	}

	fmt.Println(line)
	
	f.Close()

	return bytes, err
}

func HandleChatMessage(s *melody.Session, msg []byte) (out []byte, err error) {
	
	// Timestamp.
	timeUnixNano := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	timeString := time.Now().UTC().String()

	// IP
	ip, err := lib.GetClientIPHelper(s.Request)
	if err != nil {
		log.Fatalln("Error getting client IP: ", err)
	}

	// Geo from IP. 
	geoip, err := lib.GetGeoFromIP(ip)
	if err != nil {
		log.Fatalln("Error getting Geo IP.", err)
	}

	// JSON objectify. 
	data := j.Object().Put("time", timeString). // btw hanging dots are no go
							Put("unixNano", timeUnixNano).
							Put("message", string(msg)).
							Put("ip", ip).
							Put("bootsIP", lib.BootsEncoded(ip)).
							Put("lat", geoip["lat"]).
							Put("lon", geoip["lon"]).
							Put("city", geoip["city"]).
							Put("subdiv", geoip["subdiv"]).
							Put("countryIsoCode", geoip["countryIsoCode"]).
							Put("tz", geoip["tz"])

	dataIndentedString := data.String()
	out = []byte(dataIndentedString)

	bytes, err := saveChat(out)
	fmt.Printf("Wrote %d bytes to file\n", bytes)
	
	if err != nil {
		return nil, err
	}
	return out, nil
}