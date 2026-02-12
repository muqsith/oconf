package main

import (
	"fmt"
	"log"
	"os"

	"github.com/muqsith/oconf"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide valid config file path")
	}

	configFilePath := os.Args[1]
	jsonString, err := oconf.GetConfigAsJSONString(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println(jsonString)
}
