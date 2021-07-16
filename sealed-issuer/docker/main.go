package main

// TODO: improve error handling everywhere

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port = 8080

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/api/v1/seal", sealHandler)

	log.Printf("Listening on port %v!\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err.Error())
	}
}

func sealHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// decode
		decoder := json.NewDecoder(r.Body)
		var b SealedRequest
		err := decoder.Decode(&b)
		// response
		responseBody := SealedResponse{}
		if err != nil {
			responseBody.ErrorMessage = err.Error()
		} else {
			responseBody = Seal(b)
		}
		js, err := json.Marshal(responseBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	http.NotFound(w, r)
}
