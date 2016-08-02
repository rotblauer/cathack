package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"../lib"
	"github.com/boltdb/bolt"
	"github.com/olahol/melody"
	ghfmd "github.com/shurcooL/github_flavored_markdown"
	"sort"
)

// ChatMessageForm is exported.
type ChatMessageForm struct {
	ID             string `json:"id"`
	Time           string `json:"time"`
	UnixNano       string `json:"unixnano"`
	Message        string `json:"message"`
	IP             string `json:"ip"`
	BootsIP        string `json:"bootsIp"`
	Lat            string `json:"lat"`
	Lon            string `json:"lon"`
	City           string `json:"city"`
	Subdiv         string `json:"subdiv"`
	CountryIsoCode string `json:"countryIsoCode"`
	Tz             string `json:"tz"`
}

type ChatMessages []ChatMessageForm

func (cm *ChatMessageForm) setTimeStamp() {
	cm.Time = time.Now().UTC().String()
	cm.UnixNano = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}
func (cm *ChatMessageForm) setChatID() {
	cm.ID = lib.RandSeq(20)
}

// the ChatMessage type will implement all the methods to satisfy
// the sort.Interface interface
func (slice ChatMessages) Len() int {
	return len(slice)
}
func (slice ChatMessages) Less(i, j int) bool {
	ii, _ := strconv.Atoi(slice[i].UnixNano)
	jj, _ := strconv.Atoi(slice[j].UnixNano)
	return ii < jj
}
func (slice ChatMessages) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// AllChatMsgs does query for all chat messages in bucket "chat."
func AllChatMsgs() (ChatMessages, error) {

	var msgs ChatMessages
	var err error

	err = GetDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("chat"))

		if b.Stats().KeyN > 0 {

			c := b.Cursor()
			for chatkey, chatval := c.First(); chatkey != nil; chatkey, chatval = c.Next() {
				var msg ChatMessageForm
				json.Unmarshal(chatval, &msg)
				msgs = append(msgs, msg)
			}
			return nil // return nil (no error, and msgs are gotten)
		}
		return nil
	})
	sort.Sort(msgs)
	return msgs, err
}

// SaveChatMsg processes and formats incoming chat messages.
func SaveChatMsg(s *melody.Session, msg []byte) ([]byte, error) {

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

	markdowned := ghfmd.Markdown(msg)

	// From message struct.
	newChatMessage := ChatMessageForm{
		Message:        string(markdowned),
		IP:             ip,
		BootsIP:        lib.BootsEncoded(ip),
		Lat:            geoip["lat"],
		Lon:            geoip["lon"],
		City:           geoip["city"],
		Subdiv:         geoip["subdiv"],
		CountryIsoCode: geoip["countryIsoCode"],
		Tz:             geoip["tz"],
	}

	newChatMessage.setTimeStamp()
	newChatMessage.setChatID()

	chatMsgJSON, err := json.Marshal(newChatMessage) // return bytes, err

	if err != nil {
		fmt.Println(err)
	}

	// This can go in a go routine --
	go func() {
		GetDB().Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("chat"))
			e := b.Put([]byte(newChatMessage.ID), chatMsgJSON)
			if e != nil {
				fmt.Println(e)
				return e
			}
			return nil
		})
	}()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return chatMsgJSON, nil
}
