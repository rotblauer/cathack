package chatty

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/olahol/melody"

	"../lib"
)

// http://stackoverflow.com/questions/26327391/go-json-marshalstruct-returns
type ChatMessageAs struct {
	Time           string
	UnixNano       string
	Message        string
	Ip             string
	BootsIP        string
	Lat            string
	Lon            string
	City           string
	Subdiv         string
	CountryIsoCode string
	Tz             string
}

func saveChat(data []byte) (bytes int, err error) {

	fmt.Printf("Saving chat: %v", string(data))
	fmt.Println()

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

func HandleChatMessage(s *melody.Session, msg []byte) ([]byte, error) {

	fmt.Printf("Handling Chat Message: %v", string(msg))
	fmt.Println()

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
		Time:           timeString,
		UnixNano:       timeUnixNano,
		Message:        string(msg),
		Ip:             ip,
		BootsIP:        lib.BootsEncoded(ip),
		Lat:            geoip["lat"],
		Lon:            geoip["lon"],
		City:           geoip["city"],
		Subdiv:         geoip["subdiv"],
		CountryIsoCode: geoip["countryIsoCode"],
		Tz:             geoip["tz"], // hanging comma necessary?! https://golang.org/pkg/encoding/json/#example_Marshal
	}

	fmt.Printf("newChatMessage constructed: %v", newChatMessage)
	fmt.Println()

	// func Marshal(v interface{}) ([]byte, error)
	out, err := json.Marshal(newChatMessage)

	fmt.Printf("Marshaled out: %v", string(out))
	fmt.Println()

	if err != nil {
		fmt.Printf("Error marshaling newChatMessage: %v", err)
		fmt.Println()
	}

	bytes, err := saveChat(out)
	fmt.Printf("Wrote %d bytes to file\n", bytes)
	fmt.Println()

	if err != nil {
		fmt.Printf("There was an error in HandleChatMessage. %v", err.Error())
		fmt.Println()
		return nil, err
	}
	return out, nil
}
