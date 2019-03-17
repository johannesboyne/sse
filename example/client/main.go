package main

import (
	"bufio"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:8080"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')

		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(line))
	}
}
