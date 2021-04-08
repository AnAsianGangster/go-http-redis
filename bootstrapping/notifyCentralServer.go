package bootstrapping

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Notify(port string) {
	requestBody, err := json.Marshal(map[string]string{
		"name":   os.Getenv("SERVER_NAME"), // FIXME make server name dynamic
		"status": "ALIVE",
	})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post("http://go-consistent-hashing:"+port+"/node-status", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}
