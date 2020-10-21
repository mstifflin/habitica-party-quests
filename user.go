package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	userURL = "https://habitica.com/api/v3/members/%s"
)

func getPartyID(userUUID string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(userURL, userUUID))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	if data["error"] != nil {
		return "", errors.New(data["error"].(string))
	}

	partyIDInterface := parseUnmarshaledArbitraryJSON(data, []string{"data", "party", "_id"})
	if partyID, ok := partyIDInterface.(string); ok {
		return partyID, nil
	}

	return "", errors.New("party id not found")
}

func getUserQuests(userUUID string) (map[string]int, error) {
	resp, err := http.Get(fmt.Sprintf(userURL, userUUID))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	if data["error"] != nil {
		return nil, errors.New(data["error"].(string))
	}

	questsInterface := parseUnmarshaledArbitraryJSON(data, []string{"data", "items", "quests"})
	questsInterfaceMap, ok := questsInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("quests is not a map[string]interface")
	}

	quests := map[string]int{}

	for questName, questCountInterface := range questsInterfaceMap {
		questCount, ok := questCountInterface.(float64)
		if !ok {
			return nil, errors.New("quest count interface is not float64")
		}

		if questCount == 0 {
			continue
		}

		quests[questName] = int(questCount)
	}

	return quests, nil
}
