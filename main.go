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
	fmt.Println("apiKey", apiKey)

	partyID, err := getPartyID(myUUID)
	if err != nil {
		fmt.Println("error getting party id", err.Error())
		return
	}

	partyMembers, err := getPartyMembers(myUUID, apiKey, partyID)
	if err != nil {
		fmt.Println("error getting party members", err.Error())
		return
	}

	for _, member := range partyMembers {
		fmt.Println(member)
	}
}

func parseUnmarshaledArbitraryJSON(data map[string]interface{}, location []string) interface{} {
	if len(location) == 0 {
		return nil
	}

	nestedData, ok := data[location[0]].(map[string]interface{})
	if !ok {
		return data[location[0]]
	}

	return parseUnmarshaledArbitraryJSON(nestedData, location[1:])
}
