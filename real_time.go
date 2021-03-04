package groupme

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/karmanyaahm/wray"
)

const (
	pushServer       = "https://push.groupme.com/faye"
	userChannel      = "/user/"
	groupChannel     = "/group/"
	dmChannel        = "/direct_message/"
	handshakeChannel = "/meta/handshake"
	connectChannel   = "/meta/connect"
	subscribeChannel = "/meta/subscribe"
)

var concur = sync.Mutex{}
var token string

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

func init() {
	wray.RegisterTransports([]wray.Transport{&wray.HTTPTransport{}})
}

type HandlerAll interface {
	Handler
	HandlerText
	HandlerLike
	HandlerMembership
	HandleGroupMembership
	HandleGroupMetadata
}
type Handler interface {
	HandleError(error)
}
type HandlerText interface {
	HandleTextMessage(Message)
}
type HandlerLike interface {
	HandleLike(messageID ID, favBy []string)
}
type HandlerMembership interface {
	HandleJoin(ID)
}

type HandleGroupMetadata interface {
	HandleGroupTopic(group ID, newTopic string)
	HandleGroupName(group ID, newName string)
	HandleGroupAvatar(group ID, newAvatar string)
	HandleLikeIcon(group ID, PackID, PackIndex int, Type string)
}

type HandleGroupMembership interface {
	HandleNewNickname(group ID, user ID, newName string)
	HandleNewAvatarInGroup(group ID, user ID, avatarURL string)
}

//PushSubscription manages real time subscription
type PushSubscription struct {
	channel       chan wray.Message
	fayeClient    *wray.FayeClient
	handlers      []Handler
	LastConnected int64
}

//NewPushSubscription creates and returns a push subscription object
func NewPushSubscription(context context.Context) PushSubscription {

	r := PushSubscription{
		channel: make(chan wray.Message),
	}

	return r
}

func (r *PushSubscription) AddHandler(h Handler) {
	r.handlers = append(r.handlers, h)
}

//AddFullHandler is the same as AddHandler except to ensure interface implements everything
func (r *PushSubscription) AddFullHandler(h HandlerAll) {
	r.handlers = append(r.handlers, h)
}

type systemMessage struct {
	Event struct {
		Kind string `json:"type"`
		Data interface{}
	}
}

//Listen connects to GroupMe. Runs in Goroutine.
func (r *PushSubscription) StartListening(context context.Context) {
	r.fayeClient = wray.NewFayeClient(pushServer)

	r.fayeClient.SetLogger(fayeLogger{})

	r.fayeClient.AddExtension(&authExtension{})
	//r.fayeClient.AddExtension(r.fayeClient) //verbose output

	go r.fayeClient.Listen()

	go func() {
		for msg := range r.channel {
			r.LastConnected = time.Now().Unix()
			data := msg.Data()
			content, _ := data["subject"]
			contentType := data["type"].(string)
			channel := msg.Channel()

			if strings.HasPrefix(channel, groupChannel) || strings.HasPrefix(channel, dmChannel) {
				r.chatEvent(contentType, content)
			}

			switch contentType {
			case "line.create":
				b, _ := json.Marshal(content)
				out := Message{}
				_ = json.Unmarshal(b, &out)

				if out.UserID.String() == "system" {
					event := systemMessage{}
					err := json.Unmarshal(b, &event)
					if err != nil {
						fmt.Println(err)
					}

					r.systemEvent(out.GroupID, event)
					break
				}

				for _, h := range r.handlers {
					if h, ok := h.(HandlerText); ok {
						h.HandleTextMessage(out)
					}
				}

				break
			case "like.create":
				//should be an associated chatEvent
				break
			case "membership.create":
				c, _ := content.(map[string]interface{})
				id, _ := c["id"].(string)

				for _, h := range r.handlers {
					if h, ok := h.(HandlerMembership); ok {
						h.HandleJoin(ID(id))
					}
				}

				break
			case "ping":
				break
			default: //TODO: see if any other types are returned
				if len(contentType) == 0 || content == nil {
					break
				}
				log.Println(contentType)
				b, _ := json.Marshal(content)
				log.Fatalln(string(b))

			}

		}
	}()
}

func (r *PushSubscription) chatEvent(contentType string, content interface{}) {
	switch contentType {
	case "favorite":
		b, ok := content.(map[string]interface{})["line"].(Message)

		if !ok {
			log.Println(content)
		}

		for _, h := range r.handlers {
			if h, ok := h.(HandlerLike); ok {
				h.HandleLike(b.UserID, b.FavoritedBy)
			}
		}
		break
	default: //TODO: see if any other types are returned
		println("HEHE")
		log.Println(contentType)
		b, _ := json.Marshal(content)
		log.Fatalln(string(b))
	}

}

func (r *PushSubscription) systemEvent(groupID ID, msg systemMessage) {
	kind := msg.Event.Kind
	b, _ := json.Marshal(msg.Event.Data)
	switch kind {
	case "membership.nickname_changed":
		data := struct {
			Name string
			User struct {
				ID int
			}
		}{}
		_ = json.Unmarshal(b, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupMembership); ok {
				h.HandleNewNickname(groupID, ID(strconv.Itoa(data.User.ID)), data.Name)
			}
		}
		break
	case "membership.avatar_changed":
		data := struct {
			AvatarURL string `json:"avatar_url"`
			User      struct {
				ID int
			}
		}{}
		_ = json.Unmarshal(b, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupMembership); ok {
				h.HandleNewAvatarInGroup(groupID, ID(strconv.Itoa(data.User.ID)), data.AvatarURL)
			}
		}
		break
	case "group.name_change":
		data := struct {
			Name string
		}{}
		_ = json.Unmarshal(b, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupMetadata); ok {
				h.HandleGroupName(groupID, data.Name)
			}
		}
		break
	case "group.topic_change":
		data := struct {
			Topic string
		}{}
		_ = json.Unmarshal(b, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupMetadata); ok {
				h.HandleGroupTopic(groupID, data.Topic)
			}
		}
		break
	case "group.avatar_change":
		data := struct {
			AvatarURL string `json:"avatar_url"`
		}{}
		_ = json.Unmarshal(b, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupMetadata); ok {
				h.HandleGroupAvatar(groupID, data.AvatarURL)
			}
		}
		break
	case "group.like_icon_set":
		data := struct {
			LikeIcon struct {
				PackID    int `json:"pack_id"`
				PackIndex int `json:"pack_index"`
				Type      string
			} `json:"like_icon"`
		}{}
		_ = json.Unmarshal(b, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupMetadata); ok {
				h.HandleLikeIcon(groupID, data.LikeIcon.PackID, data.LikeIcon.PackIndex, data.LikeIcon.Type)
			}
		}
		break
	case "group.like_icon_removed":
		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupMetadata); ok {
				h.HandleLikeIcon(groupID, 0, 0, "")
			}
		}
		break
	default:
		log.Println(kind)
		log.Fatalln(string(b))
	}
}

//SubscribeToUser to users
func (r *PushSubscription) SubscribeToUser(context context.Context, userID ID, authToken string) error {
	concur.Lock()
	defer concur.Unlock()

	if r.fayeClient == nil {
		return errors.New("Not Listening") //TODO: Proper error
	}

	token = authToken
	r.fayeClient.WaitSubscribe(userChannel+userID.String(), r.channel)

	return nil
}

//SubscribeToGroup to groups for typing notification
func (r *PushSubscription) SubscribeToGroup(context context.Context, groupID ID, authToken string) error {
	concur.Lock()
	defer concur.Unlock()
	if r.fayeClient == nil {
		return errors.New("Not Listening") //TODO: Proper error
	}

	token = authToken
	r.fayeClient.WaitSubscribe(groupChannel+groupID.String(), r.channel)

	return nil
}

//Connected check if connected
func (r *PushSubscription) Connected() bool {
	return r.LastConnected+30 >= time.Now().Unix()
}

// Stop listening to GroupMe after completing all other actions scheduled first
func (r *PushSubscription) Stop(context context.Context) {
	concur.Lock()
	defer concur.Unlock()

	//TODO: stop listening
}

type authExtension struct {
}

// In does nothing in this extension, but is needed to satisy the interface
func (e *authExtension) In(msg wray.Message) {
	println(msg.Channel())
	if len(msg.Error()) > 0 {
		log.Fatalln(msg.Error())
	}
}

// Out adds the authentication token to the messages ext field
func (e *authExtension) Out(msg wray.Message) {
	if msg.Channel() == subscribeChannel {
		ext := msg.Ext()
		ext["access_token"] = token
		ext["timestamp"] = time.Now().Unix()
	}
}
