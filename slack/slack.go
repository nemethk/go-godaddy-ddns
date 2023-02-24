package slack

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func SlackPost(BEARER_TOKEN string, CHANNEL_ID string, MESSAGE string) {

    // json
    type Text struct {
        Type string `json:"type"`
        Text string `json:"text"`
    }
    type Blocks struct {
        Type string `json:"type"`
        Text Text   `json:"text"`
    }
    type Data struct {
        Channel string   `json:"channel"`
        Blocks  []Blocks `json:"blocks"`
    }

    data := Data{
        Channel: CHANNEL_ID,
        Blocks: []Blocks{
            {
                Type: "section",
                Text: Text{"mrkdwn", MESSAGE},
            },
        },
    }

    jsonData, _ := json.Marshal(data)

    // POST
    url := "https://slack.com/api/chat.postMessage"
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Add("Content-Type", "application/json; charset=utf-8")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", BEARER_TOKEN))
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    println(string(body))
}
