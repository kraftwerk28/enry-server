package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-enry/go-enry/v2"
)

var availableLanguages []string = []string{
	"Python",
	"Haskell",
	"Javascript",
	"Typescript",
	"C++",
	"C",
	"Go",
	"Lua",
	"Rust",
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	languages := req.URL.Query()["languages"]
	if len(languages) == 0 {
		languages = availableLanguages
	} else {
		languages = strings.Split(languages[0], ",")
	}
	if len(languages) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty language list is not allowed"))
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	lang := enry.GetLanguagesByClassifier("", body, languages)
	result := strings.ToLower(lang[0])
	w.Write([]byte(result))
}

func main() {
	http.HandleFunc("/", handleIndex)
	port := (os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
