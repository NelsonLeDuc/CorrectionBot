package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service"
)

func HandleCalman(message service.Message, service service.Service, cache cache.QueryCache) {

	if message.UserType() != "user" {
		return
	}

	bot, _ := models.FetchBot(message.GroupID())

	if IsImgur(message.Text()) {
		fix := CorrectImgur(message.Text())
		service.PostText(bot.Key, fix, -1, message)
	}
}

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

	result := firstGood(stuff.Results())

	return result
}

type ImgurResponse struct {
	Data struct {
		Link string `json:"link"`
		Gifv string `json:"gifv"`
		Mp4  string `json:"mp4"`
		Webm string `json:"webm"`
	} `json:"data"`
}

func (i ImgurResponse) Results() []string {
	return []string{
		i.Data.Webm,
		i.Data.Mp4,
		i.Data.Gifv,
		i.Data.Link,
	}
}

func firstGood(s []string) string {
	for _, x := range s {
		if len(x) > 0 {
			return x
		}
	}

	return ""
}
