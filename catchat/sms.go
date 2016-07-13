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

// -------------------------------------------------
// // sms("12182606849", "DDF", "c330fe3b", "d69e9ca6c8245f6a")

// //SMS text sender, nexmo to test...need a sign up with keys
// func sms(number string, messageToSend string, key string, secret string) {
// 	nexmoClient, _ := nexmo.NewClientFromAPI(key, secret)
// 	// https://github.com/njern/gonexmo
// 	// Send an SMS
// 	// See https://docs.nexmo.com/index.php/sms-api/send-message for details.
// 	message := &nexmo.SMSMessage{
// 		From:            "12529178592",
// 		To:              number,
// 		Type:            nexmo.Text,
// 		Text:            messageToSend,
// 		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
// 		Class:           nexmo.Standard,
// 	}

// 	messageResponse, err := nexmoClient.SMS.Send(message)
// 	if err != nil {
// 		log.Fatalln("Error getting sending sms: ", err)
// 	}
// 	fmt.Println(messageResponse)
// }
// -------------------------------------------------
//
// // There is https://github.com/dietsche/textbelt, but sticking with
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
