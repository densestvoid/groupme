package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/densestvoid/groupme"
)

// This is not a real token. Please find yours by logging
// into the GroupMe development website: https://dev.groupme.com/

//var authorizationToken = "ABCD"

var authorizationToken = "aa608b00a46401385ead62dd938575cf"

// A short program that gets the gets the first 5 groups
// the user is part of, and then the first 10 messages of
// the first group in that list
func main() {
	// Create a new client with your auth token
	client := groupme.NewClient(authorizationToken)

	// Get the groups your user is part of
	groups, err := client.IndexGroups(
		context.Background(),
		&groupme.GroupsQuery{
			Page:    0,
			PerPage: 1,
			Omit:    "memberships",
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Get first 10 messages of the first group
	if len(groups) == 0 {
		fmt.Println("No groups")
		os.Exit(1)
	}

	p := groupme.NewPushSubscription(context.Background())
	go p.StartListening(context.TODO())

	client = groupme.NewClient(authorizationToken)

	User, _ := client.MyUser(context.Background())
	err = p.SubscribeToUser(context.Background(), User.ID, authorizationToken)
	if err != nil {
		log.Fatal(err)
	}

	for _, j := range groups {
		err = p.SubscribeToGroup(context.TODO(), j.ID, authorizationToken)
		if err != nil {
			log.Fatal(err)
		}
	}

	p.AddFullHandler(Handler{User: User})

	<-make(chan (struct{}))
}

type Handler struct {
	User *groupme.User
}

func (h Handler) HandleError(e error) {
	fmt.Println(e)
}

func (h Handler) HandleTextMessage(msg groupme.Message) {
	fmt.Println(msg.Text, msg.Name, msg.Attachments)
}

func (h Handler) HandleJoin(group groupme.ID) {
	fmt.Println("User joined group with id", group.String())
}

func (h Handler) HandleLike(id groupme.ID, by []string) {
	fmt.Println(id.String(), "liked by", by)
}

func (h Handler) HandlerMembership(i groupme.ID) {
	fmt.Println("Membership event on", i.String())
}

func (h Handler) HandleGroupTopic(group groupme.ID, newTopic string) {
	fmt.Println(group.String(), "has new topic of", newTopic)
}
func (h Handler) HandleGroupName(group groupme.ID, newName string) {
	fmt.Println(group.String(), "has new name of", newName)
}
func (h Handler) HandleGroupAvatar(group groupme.ID, newAvatar string) {
	fmt.Println(group.String(), "has new avatar url of", newAvatar)
}

func (h Handler) HandleLikeIcon(group groupme.ID, PackID, PackIndex int, Type string) {
	//Not sure how to use without groupme icon packs
	if len(Type) == 0 {
		fmt.Println("Default like icon set")
		return
	}
	fmt.Println(group.String(), "has new like icon of", PackID, PackIndex, Type)
}

func (h Handler) HandleNewNickname(group groupme.ID, user groupme.ID, newName string) {
	fmt.Printf("In group %s, user %s has new nickname %s\n", group.String(), user.String(), newName)
}
func (h Handler) HandleNewAvatarInGroup(group groupme.ID, user groupme.ID, avatarURL string) {
	if avatarURL == "" {
		//get default avatar
		avatarURL = h.User.ImageURL
	}
	fmt.Printf("In group %s, user %s has new avatar with url %s\n", group.String(), user.String(), avatarURL)
}

func (h Handler) HandleMembers(group groupme.ID, members []groupme.Member, added bool) {
	action := "removed"
	if added {
		action = "added"
	}

	fmt.Printf("In group %s, users %v %s\n", group.String(), members, action)
}
