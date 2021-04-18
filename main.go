package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/csv"
	"io"
)

// Item representation
type Item struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func singleCsvHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Allowed POST method only", http.StatusMethodNotAllowed)
	}

	err := r.ParseMultipartForm(32 << 20) //max memory
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	file, header, err := r.FormFile("csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer file.Close()

	f, err := os.Open("file.csv")
	if err != nil {
		log.Fatal(err)
	}

	read := csv.NewReader(f)

	for {
		record, err := read.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}
	log.Fatal(header)
}

func main() {
	http.HandleFunc("/singleImport", singleCsvHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("error starting http server: ", err)
	}
}