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

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("POST /ascii", Ascii)
	http.HandleFunc("GET /ascii", AsciiMethodNotAllowed) // 405
	http.HandleFunc("/", Home)
	fmt.Println("Listening on port 8080... http://localhost:8080/")
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
	fmt.Println("raw body: ", string(body))

	incoming := myStruct{}
	err := json.Unmarshal(body, &incoming)
	if err != nil {
		fmt.Printf("\033[31mUnmarshal error: %v\033[0m\n", err)
		// Handle 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/html") // Make sure the content type is set
		http.ServeFile(w, r, "./templates/400/400.html")
		return
	}

	// check if the struct is valid
	if incoming.TypeOfAscii == "" || incoming.MakeAscii == "" {
		// Handle 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./templates/400/400.html")
		return
	}

	if incoming.TypeOfAscii != "standard" && incoming.TypeOfAscii != "shadow" && incoming.TypeOfAscii != "thinkertoy" {
		// Handle 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "./templates/404/404.html")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if !checkValidChars(incoming.MakeAscii) {
		incoming.MakeAscii = ""
	}
	asciiResponse, err := giveAscii(incoming.MakeAscii, incoming.TypeOfAscii)
	if err != nil {
		// Handle 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "./templates/500/500.html")
		return
	}

	// convert asciiResponse to json
	jsonResponse, err := json.Marshal(asciiResponse)
	if err != nil {
		// Handle 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "./templates/500/500.html")
		return
	}

	// Send the final response
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

func giveAscii(myArgs, typeOfAscii string) (string, error) {
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
