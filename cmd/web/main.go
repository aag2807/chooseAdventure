package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"makeAdventure"
)

func main() {
	port := flag.Int("port", 7070, "the default port")
	fileName := flag.String("File", "gopher.json", "The json file with the story")
	flag.Parse()
	f, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	story, err := makeAdventure.JsonStory(f)
	if err != nil {
		panic(err)
	}

	handler := makeAdventure.NewHandler(story)
	fmt.Printf("running in port %d \n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
