package lib

import (
	"strings"
	"strconv"
)

func BootsEncoded(ip string) string {
	code := strings.Split("bootsiscat", "")
	out := ""
	stringWithJustNumbers := strings.Replace(ip, ".", "0", -1)
	splitString := strings.Split(stringWithJustNumbers, "")
	for i := 0; i < len(splitString); i += 1 {
		letter := splitString[i]
		// fmt.Println(reflect.TypeOf(letter))
		number, err := strconv.Atoi(letter)
		if err != nil {
			// fmt.Println("invalid string", err)
		}
		codedLetter := code[number]
		out += codedLetter
	}
	return out
	// fmt.Println(out)
}