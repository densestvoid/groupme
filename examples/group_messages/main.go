package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/densestvoid/groupme"
)

// This is not a real token. Please find yours by logging
// into the GroupMe development website: https://dev.groupme.com/
const authorizationToken = "0123456789ABCDEF"

// A short program that gets the gets the first 5 groups
// the user is part of, and then the first 10 messages of
// the first group in that list
func main() {
	// Create a new client with your auth token
	client := groupme.NewClient(
		authorizationToken,
		groupme.WithLogHander(slog.NewJSONHandler(os.Stdout, nil)),
	)

	// Get the groups your user is part of
	groups, err := client.IndexGroups(
		context.Background(),
		&groupme.GroupsQuery{
			Page:    0,
			PerPage: 5,
			Omit:    "memberships",
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(groups)

	// Get first 10 messages of the first group
	if len(groups) == 0 {
		fmt.Println("No groups")
	}

	messages, err := client.IndexMessages(context.Background(), groups[0].ID, &groupme.IndexMessagesQuery{
		Limit: 10,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(messages)
}
