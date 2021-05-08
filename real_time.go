package groupme

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/karmanyaahm/wray"
)

const (
	PushServer       = "https://push.groupme.com/faye"
	userChannel      = "/user/"
	groupChannel     = "/group/"
	dmChannel        = "/direct_message/"
	subscribeChannel = "/meta/subscribe"
)

var (
	ErrHandlerNotFound    = errors.New("Handler not found")
	ErrListenerNotStarted = errors.New("GroupMe listener not started")
)

var concur = sync.Mutex{}
var token string

func init() {
	wray.RegisterTransports([]wray.Transport{&wray.HTTPTransport{}})
}

type HandlerAll interface {
	Handler

	//of self
	HandlerText
	HandlerLike
	HandlerMembership

	//of group
	HandleGroupTopic
	HandleGroupAvatar
	HandleGroupName
	HandleGroupLikeIcon

	//of group members
	HandleMemberNewNickname
	HandleMemberNewAvatar
	HandleMembers
}
type Handler interface {
	HandleError(error)
}
type HandlerText interface {
	HandleTextMessage(Message)
}
type HandlerLike interface {
	HandleLike(Message)
}
type HandlerMembership interface {
	HandleJoin(ID)
}

//Group Handlers
type HandleGroupTopic interface {
	HandleGroupTopic(group ID, newTopic string)
}

type HandleGroupName interface {
	HandleGroupName(group ID, newName string)
}
type HandleGroupAvatar interface {
	HandleGroupAvatar(group ID, newAvatar string)
}
type HandleGroupLikeIcon interface {
	HandleLikeIcon(group ID, PackID, PackIndex int, Type string)
}

//Group member handlers
type HandleMemberNewNickname interface {
	HandleNewNickname(group ID, user ID, newName string)
}

type HandleMemberNewAvatar interface {
	HandleNewAvatarInGroup(group ID, user ID, avatarURL string)
}
type HandleMembers interface {
	//HandleNewMembers returns only partial member with id and nickname; added is false if removing
	HandleMembers(group ID, members []Member, added bool)
}

type PushMessage interface {
	Channel() string
	Data() map[string]interface{}
	Ext() map[string]interface{}
	Error() string
}

type FayeClient interface {
	//Listen starts a blocking listen loop
	Listen()
	//WaitSubscribe is a blocking/synchronous subscribe method
	WaitSubscribe(channel string, msgChannel chan PushMessage)
}

//PushSubscription manages real time subscription
type PushSubscription struct {
	channel       chan PushMessage
	fayeClient    FayeClient
	handlers      []Handler
	LastConnected int64
}

//NewPushSubscription creates and returns a push subscription object
func NewPushSubscription(context context.Context) PushSubscription {

	r := PushSubscription{
		channel: make(chan PushMessage),
	}

	return r
}

func (r *PushSubscription) AddHandler(h Handler) {
	r.handlers = append(r.handlers, h)
}

//AddFullHandler is the same as AddHandler except it ensures the interface implements everything
func (r *PushSubscription) AddFullHandler(h HandlerAll) {
	r.handlers = append(r.handlers, h)
}

var RealTimeHandlers map[string]func(r *PushSubscription, channel string, data ...interface{})
var RealTimeSystemHandlers map[string]func(r *PushSubscription, channel string, id ID, rawData []byte)

//Listen connects to GroupMe. Runs in Goroutine.
func (r *PushSubscription) StartListening(context context.Context, client FayeClient) {
	r.fayeClient = client

	go r.fayeClient.Listen()

	go func() {
		for msg := range r.channel {
			r.LastConnected = time.Now().Unix()
			data := msg.Data()
			content := data["subject"]
			contentType := data["type"].(string)
			channel := msg.Channel()

			handler, ok := RealTimeHandlers[contentType]
			if !ok {
				if contentType == "ping" ||
					len(contentType) == 0 ||
					content == nil {
					continue
				}
				log.Println("Unable to handle GroupMe message type", contentType)
			}

			handler(r, channel, content)
		}
	}()
}

//SubscribeToUser to users
func (r *PushSubscription) SubscribeToUser(context context.Context, id ID, authToken string) error {
	return r.subscribeWithPrefix(userChannel, context, id, authToken)
}

//SubscribeToGroup to groups for typing notification
func (r *PushSubscription) SubscribeToGroup(context context.Context, id ID, authToken string) error {
	return r.subscribeWithPrefix(groupChannel, context, id, authToken)
}

//SubscribeToDM to users
func (r *PushSubscription) SubscribeToDM(context context.Context, id ID, authToken string) error {
	id = ID(strings.Replace(id.String(), "+", "_", 1))
	return r.subscribeWithPrefix(dmChannel, context, id, authToken)
}

func (r *PushSubscription) subscribeWithPrefix(prefix string, context context.Context, groupID ID, authToken string) error {
	concur.Lock()
	defer concur.Unlock()
	if r.fayeClient == nil {
		return ErrListenerNotStarted
	}

	token = authToken
	r.fayeClient.WaitSubscribe(prefix+groupID.String(), r.channel)

	return nil
}

//Connected check if connected
func (r *PushSubscription) Connected() bool {
	return r.LastConnected+30 >= time.Now().Unix()
}

// Out adds the authentication token to the messages ext field
func OutMsgProc(msg PushMessage) {
	if msg.Channel() == subscribeChannel {
		ext := msg.Ext()
		ext["access_token"] = token
		ext["timestamp"] = time.Now().Unix()
	}
}
