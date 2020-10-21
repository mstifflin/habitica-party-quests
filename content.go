package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func writeUserQuestDataToMarkdown(userName string, userQuestData string) error {
	d1 := []byte(userQuestData)
	return ioutil.WriteFile(fmt.Sprintf("../hackintasks.github.io/%s.md", userName), d1, 0644)
}

func writePartyQuestDataToMarkdown(partyQuestData string) error {
	header := `# Per User Inventory
[Bravisha_Skietano](Bravisha_Skietano.md)
[Celeblessil](Celeblessil.md)
[Hed_M](Hed_M.md)
[RexRD](RexRD.md)
[Umbra_Star](Umbra_Star.md)
[WolfenEmi](WolfenEmi.md)
[becoming_tolkien](becoming_tolkien.md)
[d34dimm0rt4l](d34dimm0rt4l.md)
[dainty_moose](dainty_moose.md)
[forestwood6](forestwood6.md)
[supervxn](supervxn.md)
[ybbm](ybbm.md)

`
	d1 := []byte(header + partyQuestData)
	return ioutil.WriteFile("../hackintasks.github.io/index.md", d1, 0644)
}

func formatPartyQuestData(
	sortedQuestKeys []string,
	totalPartyQuests map[string]int,
	questToOwnersMap map[string][]string,
	questMetadata map[string]Quest,
) string {

	section := func(b *strings.Builder, title string, category string) {
		b.WriteString(fmt.Sprintf("# %s\n", title))
		for _, questID := range sortedQuestKeys {
			if questMetadata[questID].Category == category {
				b.WriteString(fmt.Sprintf(
					"### %s\n\nQuanity: %v\n\n",
					questMetadata[questID].Name, totalPartyQuests[questID],
				))

				b.WriteString("Owners: ")
				b.WriteString(strings.Join(questToOwnersMap[questID], ", "))
				b.WriteString("\n")
				b.WriteString("\n")
			}
		}
	}

	builder := strings.Builder{}

	section(&builder, "Pet Quests", "pet")
	section(&builder, "Unlockable Quests", "unlockable")
	section(&builder, "Masterclasser Quests", "gold")
	section(&builder, "Hatching Potion Quests", "hatchingPotion")
	section(&builder, "World Quests", "world")
	section(&builder, "Hourglass Quests", "timeTravelers")

	builder.WriteString("\n")

	return builder.String()
}

func formatUserQuestData(
	sortedQuestKeys []string,
	userName string,
	userQuests map[string]int,
	questMetadata map[string]Quest,
) string {
	section := func(b *strings.Builder, title string, category string) {
		b.WriteString(fmt.Sprintf("# %s\n", title))
		for _, questID := range sortedQuestKeys {
			if questMetadata[questID].Category == category && userQuests[questID] > 0 {
				b.WriteString(fmt.Sprintf(
					"### %s\n\nQuanity: %v\n\n",
					questMetadata[questID].Name, userQuests[questID],
				))
			}
		}
	}

	builder := strings.Builder{}

	section(&builder, "Pet Quests", "pet")
	section(&builder, "Unlockable Quests", "unlockable")
	section(&builder, "Masterclasser Quests", "gold")
	section(&builder, "Hatching Potion Quests", "hatchingPotion")
	section(&builder, "World Quests", "world")
	section(&builder, "Hourglass Quests", "timeTravelers")

	return builder.String()
}

const (
	partyChatURL = "https://habitica.com/api/v3/groups/%s/chat"
)

// Message too big for one chat message. lol
func sendChatMessageToParty(userUUID, apiKey, partyID, message string) error {
	client := &http.Client{}

	reqBody, err := json.Marshal(map[string]string{
		"message": message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(partyChatURL, partyID), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-client", myUUID)
	req.Header.Add("x-api-user", userUUID)
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	if data["error"] != nil {
		return errors.New(data["error"].(string))
	}

	return nil
}
