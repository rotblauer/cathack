package lib

import (
	"math/rand"
	"strings"
)

// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetLanguageModeByExtension(name string) (lang string) {
	exs := strings.Split(name, ".")
	for _, ex := range exs {
		switch {
		case ex == "js":
			lang = "javascript"
		case ex == "html":
			lang = "htmlmixed"
		case ex == "go":
			lang = "go"
		case ex == "md":
			lang = "markdown"
		case ex == "mdown":
			lang = "markdown"
		case ex == "markdown":
			lang = "markdown"
		case ex == "py":
			lang = "python"
		case ex == "rb":
			lang = "ruby"
		case ex == "r":
			lang = "r"
		case ex == "sh":
			lang = "shell"
		case ex == "zsh":
			lang = "shell"
		case ex == "swift":
			lang = "swift"
		default:
			lang = "markdown"
		}
	}
	return lang
}
