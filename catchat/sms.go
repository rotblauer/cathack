package catchat

import (
	"log"
	"regexp"
	"strings"

	"gopkg.in/dietsche/textbelt.v1"
)

func sendSMS(number string, message string) error {
	texter := textbelt.New()
	err := texter.Text(number, message)
	return err
}

func DelegateSendSMS(messageText []byte) (status []byte, err error) {

	// Status
	// NOMATCH - no mention of @name in the message
	// NOFIND - @name not found in phonebook
	// SENDERR - error sending <textbelt> sms
	// SENDOK - success sending <textbelt> sms

	status = []byte("NOMATCH")
	messageString := string(messageText)

	phoneBook := make(map[string]string)
	phoneBook["john"] = "2182606849"
	phoneBook["isaac"] = "2183494908"
	phoneBook["sharif"] = "6073420398"
	phoneBook["aaron"] = "3603393496"
	phoneBook["kelsey"] = "3603393496"
	phoneBook["alexa"] = "3603393496"

	re, err := regexp.Compile(`@(\w+)`) // FIXME: this should capture only the name, not the @ part. it captures @name. don't know why.
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		log.Println()
	}

	if re.MatchString(messageString) {

		status = []byte("NOFIND")

		log.Printf("Regex matches. Sending smss.")
		log.Println()

		// get @names
		names := re.FindAllString(messageString, 3) // limit to first 3 matches (from left -> right)

		// Remove @name's if we want.
		// messageString = re.ReplaceAllString(messageString, "")

		// send smss
		for _, name := range names {

			phoneNumber := phoneBook[strings.Replace(string(name), "@", "", 1)]

			log.Printf("Name: %v, Number: %v", name, phoneNumber)
			log.Println()

			// if the phone number exists for the given name
			if len(phoneNumber) > 0 {

				err := sendSMS(phoneNumber, messageString)

				if err != nil {
					log.Printf("Textbelt error: %v", err)
					log.Println()
					status = nil
				} else {
					log.Printf("Sent SMS to %v with content: %v", phoneNumber, messageString)
					// m.Broadcast([]byte("{\"status\": \"" + "Whoosh!" + "\"}"))
					status = []byte("SENDOK")
				}
			}
		}
	}

	return status, err
}
