package main

import (
	"fmt"
	"os"
	"time"
)

const (
	myUUID = "c3c09562-c4ff-4852-a69d-a9150eaf4ebc"
)

func main() {
	apiKey := os.Args[1]

	for true {
		// Get the user's party ID
		partyID, err := getPartyID(myUUID)
		if err != nil {
			fmt.Println("error getting party id", err.Error())
			return
		}

		// Given a partyID, get all the party's members names and UUIDs
		partyMembers, err := getPartyMembers(myUUID, apiKey, partyID)
		if err != nil {
			fmt.Println("error getting party members", err.Error())
			return
		}

		// Given a list of party members, get all the quests they
		// have in their inventory
		allUserQuests := []map[string]int{}
		partyMemberToQuestInventoryMap := map[string]map[string]int{}
		for _, member := range partyMembers {
			userQuests, err := getUserQuests(member.UserUUID)
			if err != nil {
				fmt.Println("error getting party members quests", err.Error())
				return
			}

			partyMemberToQuestInventoryMap[member.UserName] = userQuests
			allUserQuests = append(allUserQuests, userQuests)
		}

		// Get quest ownership list
		questToOwnersMap := questsToListOfOwnersMap(partyMemberToQuestInventoryMap)

		// Merge all the quest counts together
		sortedQuestKeys, totalPartyQuests := countTotalUserQuests(allUserQuests)

		// Get quest metadata
		questMetadata, err := getQuestData()
		if err != nil {
			fmt.Println("error getting party members", err.Error())
			return
		}

		partyQuestsString := formatPartyQuestData(
			sortedQuestKeys,
			totalPartyQuests,
			questToOwnersMap,
			questMetadata,
		)

		partyMemberNames := []string{}
		for memberName := range partyMemberToQuestInventoryMap {
			partyMemberNames = append(partyMemberNames, memberName)
		}

		err = writePartyQuestDataToMarkdown(partyMemberNames, partyQuestsString)
		if err != nil {
			fmt.Println("error writing party quest data to md file", err.Error())
			return
		}

		for userName, userQuests := range partyMemberToQuestInventoryMap {
			userQuestData := formatUserQuestData(sortedQuestKeys, userName, userQuests, questMetadata)
			err = writeUserQuestDataToMarkdown(userName, userQuestData)
			if err != nil {
				fmt.Println("error writing user quest data to md file", err.Error())
				return
			}
		}

		time.Sleep(60 * 60 * time.Second)
	}

	return
}

func parseUnmarshaledArbitraryJSON(data map[string]interface{}, location []string) interface{} {
	if len(location) == 0 {
		return nil
	}

	if len(location) == 1 {
		return data[location[0]]
	}

	nestedData, ok := data[location[0]].(map[string]interface{})
	if !ok {
		return data[location[0]]
	}

	return parseUnmarshaledArbitraryJSON(nestedData, location[1:])
}
