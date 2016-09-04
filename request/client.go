package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mikechinaloy/vindinium-go/model"
)

// PostRequest sends an HTTP request
func PostRequest(url string, body map[string]string) *model.GameState {
	var respGameState model.GameState
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshal(body)))
	if err != nil {
		fmt.Printf("\nUnable to POST request=%s", err)
	} else {
		return unmarshal(resp)
	}
	return &respGameState
}

func marshal(body map[string]string) []byte {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("\nUnable to marshal body%s", err)
	}
	return jsonBody
}

func unmarshal(resp *http.Response) *model.GameState {
	var gameState model.GameState
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\nUnable to read response=%s", err)
	} else {
		err := json.Unmarshal(body, &gameState)
		if err != nil {
			fmt.Printf("\nUnable to unmarshal body=%s", err)
		}
	}
	return &gameState
}
