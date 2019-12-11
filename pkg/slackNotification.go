package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//Payloads ...
type Payloads struct {
	Attachments map[string][]Payload `json: "attachments"`
}

//Payload ...
type Payload struct {
	Title      string `json:"title"`
	Pretext    string `json:"pretext"`
	Text       string `json:"text"`
	Color      string `json:"color"`
	AuthorName string `json:"author_name"`
	TitleLink  string `json:"title_link"`
	Footer     string `json:"footer"`
	Timestamp  string `json:"ts"`
}

//SlackNotification ...
func (p Payloads) SlackNotification(url string) {
	jsonValue, err := json.Marshal(p.Attachments)
	if err != nil {
		log.Info(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Info("response Status:", resp.Status)
}
