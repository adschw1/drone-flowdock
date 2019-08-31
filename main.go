package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var logFatalf = log.Fatalf

type flowMessage struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

type inboxMessage struct {
	Event string `json:"event"`
	Title string `json:"title"`
}

func postMessage(raw []byte, err error, flowURL string) {
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(flowURL, "application/json", bytes.NewReader(raw))
	if resp != nil {
		if resp.StatusCode != 202 {
			logFatalf("Failed to post message, flowdock api returned: %s", resp.Status)
		}
		fmt.Println(resp.Status)
		resp.Body.Close()
	}

	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	apiURL := "https://api.flowdock.com/messages?flow_token="

	flowToken := os.Getenv("PLUGIN_FLOW_TOKEN")
	if flowToken == "" {
		log.Fatalln("Missing setting: flow_token")
	}
	flowURL := apiURL + flowToken

	message := os.Getenv("PLUGIN_MESSAGE")
	if message == "" {
		repoName := os.Getenv("DRONE_REPO")
		buildLink := os.Getenv("DRONE_BUILD_LINK")
		buildStatus := os.Getenv("DRONE_BUILD_STATUS")
		message = fmt.Sprintf("Status of build [%s](%s) is %s", repoName, buildLink, buildStatus)
	}

	messageType := os.Getenv("PLUGIN_MESSAGE_TYPE")

	if messageType == "activity" {
		msg := inboxMessage{
			Event: "activity",
			Title: message,
		}

		raw, err := json.Marshal(msg)

		postMessage(raw, err, flowURL)
	} else {
		msg := flowMessage{
			Event:   "message",
			Content: message,
		}

		raw, err := json.Marshal(msg)

		postMessage(raw, err, flowURL)
	}
}
