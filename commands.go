package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Commands map[string]string

func readCommands(filename string) Commands {
	var commands Commands
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := ioutil.WriteFile(filename, commands.toJson(), 0640); err != nil {
			log.Fatal(err)
		}
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bytes, &commands); err != nil {
		log.Fatal(err)
	}

	return commands
}

func (commands *Commands) toJson() []byte {
	bytes, err := json.MarshalIndent(commands, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
