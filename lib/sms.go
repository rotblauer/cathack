package lib

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// There is https://github.com/dietsche/textbelt, but sticking with
// roll your own for now.
func SendSMS(number string, message string) (res *http.Response, err error) {

	// $ curl -X POST http://textbelt.com/text \
	//    -d number=5551234567 \
	//    -d "message=I sent this message for free with textbelt.com"

	client := http.Client{}

	textBeltUrl := "http://textbelt.com/text"

	form := url.Values{}
	form.Add("number", number)   //number
	form.Add("message", message) //message

	req, err := http.NewRequest("POST", textBeltUrl, strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	return res, err
}
