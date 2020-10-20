package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

func getPartyID(userUUID string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://habitica.com/api/v3/members/%s", userUUID))
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

type partyMember struct {
	UserName string
	UserUUID string
}

func getPartyMembers(userUUID, apiKey, partyID string) ([]partyMember, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://habitica.com/api/v3/groups/%s/members", partyID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-client", myUUID)
	req.Header.Add("x-api-user", userUUID)
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
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

	partyMembersInterface := parseUnmarshaledArbitraryJSON(data, []string{"data"})

	partyMembersJSONInterfaceList, ok := partyMembersInterface.([]interface{})
	if !ok {
		return nil, errors.New("party member list is not a []interface{}")
	}

	partyMembers := []partyMember{}

	for _, memberJSONInterface := range partyMembersJSONInterfaceList {
		memberJSON, ok := memberJSONInterface.(map[string]interface{})
		if !ok {
			return nil, errors.New("member json is not type map[string]interface{}")
		}

		userNameInterface := parseUnmarshaledArbitraryJSON(memberJSON, []string{"auth", "local", "username"})
		userName, ok := userNameInterface.(string)
		if !ok {
			return nil, errors.New("wrong user name type")
		}

		idInterface := parseUnmarshaledArbitraryJSON(memberJSON, []string{"id"})
		id, ok := idInterface.(string)
		if !ok {
			return nil, errors.New("wrong id type")
		}

		partyMembers = append(partyMembers, partyMember{UserName: userName, UserUUID: id})
	}

	return partyMembers, nil
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
