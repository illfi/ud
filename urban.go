package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DefinitionEntry struct {
	Id         int32  `json:"defid"`
	Word       string `json:"word"`
	Author     string `json:"author"`
	URL        string `json:"permalink"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
	ThumbsUp   int32  `json:"thumbs_up"`
	ThumbsDown int32  `json:"thumbs_down"`
}

type UrbanDictResp struct {
	Tags        []string          `json:"tags"`
	ResultType  string            `json:"result_type"`
	Definitions []DefinitionEntry `json:"list"`
	Sounds      []string          `json:"sounds"`
}

const (
	UD_API_URL string = "http://api.urbandictionary.com/v0/define?term=%s"
)

func UDLookup(query string) (UrbanDictResp, error) {
	resp, err := http.Get(fmt.Sprintf(UD_API_URL, url.QueryEscape(query)))
	var qResp UrbanDictResp
	if err != nil {
		return qResp, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return qResp, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal(content, &qResp)
	if err != nil {
		return qResp, nil
	}
	return qResp, nil
}

func main() {
	var top = flag.Int("top", 1, "How many top results to display")
	var audio = flag.Bool("audio", false, "return the link to the audio file")
	flag.Parse()

	var query string

	if len(flag.Args()) < 1 {
		fmt.Println("Missing Query.")
		os.Exit(1)
	} else if len(flag.Args()) > 1 {
		query = strings.Join(flag.Args(), " ")
	} else {
		query = flag.Arg(0)
	}

	resp, err := UDLookup(query)
	if err != nil {
		panic(err)
	}

	if *audio {
		var sounds []string
		if len(resp.Sounds) > *top {
			sounds = resp.Sounds[:*top]
		} else {
			sounds = resp.Sounds
		}
		for ind, sound := range sounds {
			fmt.Printf("[%d]: %s\n", ind+1, sound)
		}
	} else {
		var definitions []DefinitionEntry
		if len(resp.Definitions) > *top {
			definitions = resp.Definitions[:*top]
		} else {
			definitions = resp.Definitions
		}
		for ind, def := range definitions {
			fmt.Printf("[%d]: %s\nby %s\n\n", ind+1, def.Definition, def.Author)
		}
	}

}
