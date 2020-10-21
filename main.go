package main

import (
	"fmt"
	"os"
)

const (
	myUUID = "c3c09562-c4ff-4852-a69d-a9150eaf4ebc"
)

func main() {
	apiKey := os.Args[1]

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

	message := formatPartyQuestData(
		sortedQuestKeys,
		totalPartyQuests,
		questToOwnersMap,
		questMetadata,
	)

	fmt.Println(message)

	// Message too big for one chat message. lol
	// err = sendChatMessageToParty(myUUID, apiKey, partyID, message)
	// if err != nil {
	// 	fmt.Println("error sending message to party", err.Error())
	// 	return
	// }

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
