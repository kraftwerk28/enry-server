package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-enry/go-enry/v2"
	"github.com/julienschmidt/httprouter"
)

type LanguageCfg struct {
	name          string
	commentstring string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type DetectedLangResponse struct {
	Language      string `json:"language"`
	Extension     string `json:"extension"`
	Commentstring string `json:"commentstring"`
}

var availableLanguages = []LanguageCfg{
	{"Python", "# %s"},
	{"Haskell", "-- %s"},
	{"Javascript", "// %s"},
	{"Typescript", "// %s"},
	{"C++", "// %s"},
	{"C", "// %s"},
	{"Go", "// %s"},
	{"Lua", "-- %s"},
	{"Rust", "// %s"},
}

func normalizeLangList() {
	for i, cfg := range availableLanguages {
		oldName := cfg.name
		lang, ok := enry.GetLanguageByAlias(oldName)
		if !ok {
			log.Fatalf(
				"\"%s\" is no a valid language name among `availableLanguages`",
				oldName,
			)
		}
		availableLanguages[i].name = lang
	}
}

func getCommentString(lang string) string {
	for _, cfg := range availableLanguages {
		if cfg.name == lang {
			return cfg.commentstring
		}
	}
	return ""
}

// We assume, that Content-Type is already application/json
func fail(res http.ResponseWriter, code int, message string) {
	res.WriteHeader(code)
	body, _ := json.Marshal(ErrorResponse{message})
	res.Write(body)
}

func handleDetectLanguage(
	w http.ResponseWriter,
	req *http.Request,
	ps httprouter.Params,
) {
	w.Header().Add("Content-Type", "application/json")
	rawLanguagesFromQuery := req.URL.Query().Get("languages")
	var candidates []string

	if rawLanguagesFromQuery == "" {
		candidates = make([]string, len(availableLanguages))
		for i := range availableLanguages {
			candidates[i] = availableLanguages[i].name
		}
	} else {
		candidates = strings.Split(rawLanguagesFromQuery, ",")
	}

	if len(candidates) == 0 {
		fail(w, http.StatusBadRequest, "Empty language list is not allowed")
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fail(w, http.StatusInternalServerError, err.Error())
		return
	}

	lang, _ := enry.GetLanguageByClassifier(body, candidates)
	if lang == "" {
		fail(w, http.StatusNotFound, "Langauge not detected")
		return
	}

	// Ok will always be true
	ext := enry.GetLanguageExtensions(lang)
	commentString := getCommentString(lang)
	resBody, _ := json.Marshal(DetectedLangResponse{
		lang,
		ext[0][1:],
		commentString,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func setupServer() {
	router := httprouter.New()
	router.POST("/", handleDetectLanguage)
	port := (os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), router)
	log.Fatal(err)
}

func main() {
	normalizeLangList()
	setupServer()
}
