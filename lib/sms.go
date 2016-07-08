package lib

import (
	"gopkg.in/dietsche/textbelt.v1"
	"log"
	"regexp"
	"strings"
)

// There is https://github.com/dietsche/textbelt, but sticking with
// roll your own for now.
// func SendSMS(number string, message string) (*http.Response, error) {

// 	// $ curl -X POST http://textbelt.com/text \
// 	//    -d number=5551234567 \
// 	//    -d "message=I sent this message for free with textbelt.com"

// 	client := http.Client{}

// 	textBeltUrl := "http://textbelt.com/text"

// 	form := url.Values{}
// 	form.Add("number", number)   //number
// 	form.Add("message", message) //message

// 	req, err := http.NewRequest("POST", textBeltUrl, strings.NewReader(form.Encode()))
// 	req.PostForm = form
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	res, err := client.Do(req)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(res)

// 	return res, err
// }

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

	re, err := regexp.Compile(`@(\w+)`) // FIXME: this should capture only the name, not the @ part. it captures @name. don't know why.
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
	}

	if re.MatchString(messageString) {

		status = []byte("NOFIND")

		log.Printf("Regex matches. Sending smss.")

		// get @names
		names := re.FindAllString(messageString, 3) // limit to first 3 matches (from left -> right)

		// Remove @name's if we want.
		// messageString = re.ReplaceAllString(messageString, "")

		// send smss
		for _, name := range names {

			phoneNumber := phoneBook[strings.Replace(string(name), "@", "", 1)]

			log.Printf("Name: %v, Number: %v", name, phoneNumber)

			// if the phone number exists for the given name
			if len(phoneNumber) > 0 {

				err := sendSMS(phoneNumber, messageString)

				if err != nil {
					log.Printf("Textbelt error: %v", err)
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
