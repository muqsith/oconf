package main

import (
	"fmt"
	"log"
	"os"

	"github.com/muqsith/oconf"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide valid config file path")
	}
	configFilePath := os.Args[1]
	fmt.Println(oconf.GetConfigAsJSONString(configFilePath))
}