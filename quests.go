package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

const (
	contentURL = "https://habitica.com/api/v3/content"
)

// Quest ...
type Quest struct {
	ID       string
	Category string
	Name     string
}

func getQuestData() (map[string]Quest, error) {
	resp, err := http.Get(contentURL)
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

	questsJSONInterface := parseUnmarshaledArbitraryJSON(data, []string{"data", "quests"})

	questsJSON, ok := questsJSONInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("questsJSONInterface is not a map[string]interface")
	}

	quests := map[string]Quest{}
	for questID, valueInterface := range questsJSON {
		questJSONMap, ok := valueInterface.(map[string]interface{})
		if !ok {
			return nil, errors.New("quest data keyed by quest name is not a map[string]interface")
		}

		quests[questID] = Quest{
			ID:       questID,
			Category: fmt.Sprintf("%v", questJSONMap["category"]),
			Name:     fmt.Sprintf("%v", questJSONMap["text"]),
		}
	}

	return quests, nil
}

func countTotalUserQuests(allUserQuests []map[string]int) ([]string, map[string]int) {
	totalQuests := map[string]int{}

	for _, userQuests := range allUserQuests {
		for questID, count := range userQuests {
			totalQuests[questID] += count
		}
	}

	sortedQuestKeys := []string{}
	for questID := range totalQuests {
		sortedQuestKeys = append(sortedQuestKeys, questID)
	}

	sort.Strings(sortedQuestKeys)

	return sortedQuestKeys, totalQuests
}

func questsToListOfOwnersMap(partyMemberNameToOwnedQuestsMap map[string]map[string]int) map[string][]string {
	questIDToListOfOwnersMap := map[string][]string{}

	for memberName, ownedQuestMap := range partyMemberNameToOwnedQuestsMap {
		for questID := range ownedQuestMap {
			if _, ok := questIDToListOfOwnersMap[questID]; !ok {
				questIDToListOfOwnersMap[questID] = []string{}
			}

			questIDToListOfOwnersMap[questID] = append(questIDToListOfOwnersMap[questID], memberName)
		}
	}

	for _, listOfOwners := range questIDToListOfOwnersMap {
		sort.Strings(listOfOwners)
	}

	return questIDToListOfOwnersMap
}
