package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var myMap map[rune][8]string
var readyForPrint string

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("POST /ascii", Ascii)
	http.HandleFunc("GET /ascii", AsciiMethodNotAllowed) // 405
	http.HandleFunc("/", Home)
	fmt.Println("Listening on port 8080... \033[34mhttp://localhost:8080/\033[0m")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func AsciiMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "./templates/404/404.html")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "./templates/home/home.html")
}

func Ascii(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	incoming := myStruct{}
	err := json.Unmarshal(body, &incoming)
	if err != nil {
		fmt.Printf("\033[31mUnmarshal error: %v\033[0m\n", err)
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// check if the struct is valid
	if incoming.TypeOfAscii == "" || incoming.MakeAscii == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	if incoming.TypeOfAscii != "standard" && incoming.TypeOfAscii != "shadow" && incoming.TypeOfAscii != "thinkertoy" {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if !checkValidChars(incoming.MakeAscii) {
		incoming.MakeAscii = "Invalid characters"
	}
	asciiResponse, err := giveAsccii(incoming.MakeAscii, incoming.TypeOfAscii)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	jsonResponse, _ := json.Marshal(asciiResponse) // convert asciiResponse to json
	fmt.Fprintf(w, string(jsonResponse))
}

func checkValidChars(myArgs string) bool {
	for _, char := range myArgs {
		if (char < 32 || char > 126) && char != '\n' {
			return false
		}
	}
	return true
}

type myStruct struct {
	TypeOfAscii string `json:"type"`
	MakeAscii   string `json:"text"`
}

func giveAsccii(myArgs, typeOfAscii string) (string, error) {
	mybytes, err := os.ReadFile(typeOfAscii + ".txt")
	if err != nil {
		return "", err
	}

	myMap = fillMap(squeeze(strings.Split(strings.ReplaceAll(string(mybytes), "\r\n", "\n"), "\n")))
	words := strings.Split(myArgs, "\n")
	return giveInPrintFormatWords(words), nil
}

func giveInPrintFormatWords(words []string) (formatted string) {
	for _, word := range words {
		if word == "" {
			formatted += "\n"
		} else {
			formatted += (formatWord(word))
		}
	}
	return
}

func formatWord(word string) (formatted string) {
	for i := 0; i < 8; i++ {
		for _, rune := range word {
			formatted += myMap[rune][i]
		}
		formatted += "\n"
	}
	return
}

func fillMap(s []string) map[rune][8]string {
	result := make(map[rune][8]string)
	for i := 0; i < len(s); i += 8 {
		result[rune(i/8+' ')] = [8]string{s[i], s[i+1], s[i+2], s[i+3], s[i+4], s[i+5], s[i+6], s[i+7]}
	}
	return result
}

func squeeze(content []string) (result []string) {
	for _, line := range content {
		if line == "" {
			continue
		}
		result = append(result, line)
	}
	return
}
