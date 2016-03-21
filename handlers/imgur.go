package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const gifMaxSize = 1000000 * 20 //20 megabytes
var imgurMatches = matchFunc("https*:\\/\\/.*imgur.com(?:\\/gallery)*\\/(\\w+)(?:\\?r)*[^.]?(?:$|\\s)")

func CorrectImgur(text string) string {
	matched := imgurMatches(text)
	if len(matched) < 2 {
		return ""
	}

	id := matched[1]

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.imgur.com/3/image/"+id, nil)

	clientID := os.Getenv("imgur_id")
	req.Header.Set("Authorization", "Client-ID "+clientID)
	res, _ := client.Do(req)

	body, _ := ioutil.ReadAll(res.Body)
	var stuff ImgurResponse
	json.Unmarshal(body, &stuff)

	if strings.HasSuffix(stuff.Data.Link, ".gif") && stuff.Data.Size > gifMaxSize {
		return ""
	}

	return stuff.Data.Link
}

type ImgurResponse struct {
	Data struct {
		Size int    `json:"size"`
		Link string `json:"link"`
		Gifv string `json:"gifv"`
		Mp4  string `json:"mp4"`
		Webm string `json:"webm"`
	} `json:"data"`
}
