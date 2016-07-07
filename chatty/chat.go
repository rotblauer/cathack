package chatty

import (
	"os"
	"strconv"
	"time"
	"log"
	"fmt"
	"encoding/json"
	
	"github.com/olahol/melody"

	"../lib"
)

type ChatMessageAs struct {
	time string
	unixNano string
	message string
	ip string
	bootsIP string
	lat string
	lon string
	city string
	subdiv string
	countryIsoCode string
	tz string
}

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

	// From message struct.
	newChatMessage := ChatMessageAs{
		time: timeString,
		unixNano: timeUnixNano,
		message: string(msg),
		ip: ip,
		bootsIP: lib.BootsEncoded(ip),
		lat: geoip["lat"],
		lon: geoip["lon"],
		city: geoip["city"],
		subdiv: geoip["subdiv"],
		countryIsoCode: geoip["countryIsoCode"],
		tz: geoip["tz"]}

	out, _ = json.Marshal(newChatMessage)

	bytes, err := saveChat(out)
	fmt.Printf("Wrote %d bytes to file\n", bytes)
	
	if err != nil {
		return nil, err
	}
	return out, nil
}