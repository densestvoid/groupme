package main

import (
	"context"
	"fmt"
	"log"

	"github.com/densestvoid/groupme"
	"github.com/karmanyaahm/wray"
)

// This is not a real token. Please find yours by logging
// into the GroupMe development website: https://dev.groupme.com/

var authorizationToken = "ABCD"

//This adapts your faye library to an interface compatible with this library
type FayeClient struct {
	*wray.FayeClient
}

func (fc FayeClient) WaitSubscribe(channel string, msgChannel chan groupme.PushMessage) {
	c_new := make(chan wray.Message)
	fc.FayeClient.WaitSubscribe(channel, c_new)
	//converting between types because channels don't support interfaces well
	go func() {
		for i := range c_new {
			msgChannel <- i
		}
	}()
}

//for authentication, specific implementation will vary based on faye library
type AuthExt struct{}

func (a *AuthExt) In(wray.Message) {}
func (a *AuthExt) Out(m wray.Message) {
	groupme.OutMsgProc(m)
}

//specific to faye library
type fayeLogger struct{}

func (l fayeLogger) Infof(f string, a ...interface{}) {
	log.Printf("[INFO]  : "+f, a...)
}
func (l fayeLogger) Errorf(f string, a ...interface{}) {
	log.Printf("[ERROR] : "+f, a...)
}
func (l fayeLogger) Debugf(f string, a ...interface{}) {
	log.Printf("[DEBUG] : "+f, a...)
}
func (l fayeLogger) Warnf(f string, a ...interface{}) {
	log.Printf("[WARN]  : "+f, a...)
}

// A short program that subscribes to 2 groups and 2 direct chats
// and prints out all recognized events in those
func main() {

	//Create and initialize fayeclient
	fc := FayeClient{wray.NewFayeClient(groupme.PushServer)}
	fc.SetLogger(fayeLogger{})
	fc.AddExtension(&AuthExt{})
	//for additional logging uncomment the following line
	//fc.AddExtension(fc.FayeClient)

	//create push subscription and start listening
	p := groupme.NewPushSubscription(context.Background())
	go p.StartListening(context.TODO(), fc)

	// Create a new client with your auth token
	client := groupme.NewClient(authorizationToken)
	User, _ := client.MyUser(context.Background())
	//Subscribe to get messages and events for the specific user
	err := p.SubscribeToUser(context.Background(), User.ID, authorizationToken)
	if err != nil {
		log.Fatal(err)
	}

	//handles (in this case prints) all messages
	p.AddFullHandler(Handler{User: User})

	// Get the groups your user is part of
	groups, err := client.IndexGroups(
		context.Background(),
		&groupme.GroupsQuery{
			Page:    0,
			PerPage: 2,
			Omit:    "memberships",
		})

	if err != nil {
		fmt.Println(err)
		return
	}
	//Subscribe to those groups
	for _, j := range groups {
		err = p.SubscribeToGroup(context.TODO(), j.ID, authorizationToken)
		if err != nil {
			log.Fatal(err)
		}
	}

	//get chats your user is part of
	chats, err := client.IndexChats(context.Background(),
		&groupme.IndexChatsQuery{
			Page:    0,
			PerPage: 2,
		})
	//subscribe to all those chats
	for _, j := range chats {
		err = p.SubscribeToDM(context.TODO(), j.LastMessage.ConversationID, authorizationToken)
		if err != nil {
			log.Fatal(err)
		}
	}

	//blocking
	<-make(chan (struct{}))
}

//Following example handlers print out all data
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

func (h Handler) HandleLike(msg groupme.Message) {
	fmt.Println(msg.ID, "liked by", msg.FavoritedBy)
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
