package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func readJson(filename string, data interface{}) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := ioutil.WriteFile(filename, toJson(data), 0640); err != nil {
			log.Fatal(err)
		}
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatal(err)
	}
}

func toJson(data interface{}) []byte {
	bytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
