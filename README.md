# Habitica Party Quest Tracker

## About
This application queries Habitica for a user's party member's quests and chats a message to the party when there has been an update.

## Dependencies
- go
- jq (for development and API exploration)

## Usage
```
$ go build
$ ./habitica-party-quests <user api key>
```

## Notes
- `x-client` header should always be the developer's UUID
- `x-api-user` and `x-api-key` should use the the values of the user interacting with the tool

```
# Get a user's party ID
curl -X GET https://habitica.com/api/v3/members/<user uuid> | jq .data.party._id

# Given a party ID, get a list of its members and its member's userUUIDs
curl -X GET https://habitica.com/api/v3/groups/<party id>/members \
  -H "Content-Type:application/json" \
  -H "x-client: <user uuid>`" \
  -H "x-api-user: <user uuid>`" \
  -H "x-api-key: <user api key>" | jq '.data | .[] | {username: .auth.local.username, id: .id}'

# Post a message to a party's chat
curl -X POST https://habitica.com/api/v3/groups/<party id>/chat \
  -H "Content-Type:application/json" \
  -H "x-client: <user uuid>" \
  -H "x-api-user: <user uuid>" \
  -H "x-api-key: <user api key>" \
  -d '{"message": "test message from curl"}'

# Get a user's quests
curl -X GET https://habitica.com/api/v3/members/<user uuid> | jq .data.items.quests

# Get all quests available in Habitica
curl https://habitica.com/api/v3/content | jq '.data.quests | .[] | {id: .key, category: .category, name: .text}' >> quests.json
```
