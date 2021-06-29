package main

import (
	"flag"
	"fmt"
	"os"

	"makeAdventure"
)

func main() {
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

	fmt.Println(story)

}
