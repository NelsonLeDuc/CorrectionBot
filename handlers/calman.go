package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service"
)

func HandleCalman(message service.Message, service service.Service, cache cache.QueryCache) {

	if message.UserType() != "user" {
		return
	}

	bot, _ := models.FetchBot(message.GroupID())

	fmt.Println(message.Text())

	if IsImgur(message.Text()) {
		fmt.Println("matched")
		fix := CorrectImgur(message.Text())
		if len(fix) == 0 {
			return
		}

		fmt.Println(fix)
		service.PostText(bot.Key, fix, -1, message)
	}
}

//IMGUR

const gifMaxSize = 1000000 * 20 //20 megabytes

func IsImgur(text string) bool {
	r, _ := regexp.Compile("(?i)https*:\\/\\/w{0,3}.*imgur.com\\/gallery\\/(\\w+)")
	matched := r.FindStringSubmatch(text)
	return len(matched) >= 2
}

func CorrectImgur(text string) string {
	r, _ := regexp.Compile("(?i)https*:\\/\\/w{0,3}.*imgur.com\\/gallery\\/(\\w+)")
	matched := r.FindStringSubmatch(text)
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
