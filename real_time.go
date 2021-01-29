package groupme

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/karmanyaahm/wray"
)

const (
	pushServer       = "https://push.groupme.com/faye"
	userChannel      = "/user/"
	groupChannel     = "/group/"
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
	//	log.Printf("[DEBUG] : "+f, a...)
}
func (l fayeLogger) Warnf(f string, a ...interface{}) {
	log.Printf("[WARN]  : "+f, a...)
}

func init() {
	wray.RegisterTransports([]wray.Transport{&wray.HTTPTransport{}})
}

//LikeEvent returns events as they happen from GroupMe
type LikeEvent struct {
	Message Message
}

//PushSubscription manages real time subscription
type PushSubscription struct {
	channel chan wray.Message
	//Channel to return things
	MessageChannel chan Message
	LikeChannel    chan LikeEvent
	fayeClient     *wray.FayeClient
}

//NewPushSubscription creates and returns a push subscription object
func NewPushSubscription(context context.Context) PushSubscription {

	r := PushSubscription{
		channel:        make(chan wray.Message),
		MessageChannel: make(chan Message),
		LikeChannel:    make(chan LikeEvent),
	}

	return r
}

//Listen connects to GroupMe. Runs in Goroutine.
func (r *PushSubscription) Listen(context context.Context) {
	r.fayeClient = wray.NewFayeClient(pushServer)

	r.fayeClient.SetLogger(fayeLogger{})

	r.fayeClient.AddExtension(&authExtension{})
	//r.fayeClient.AddExtension(r.fayeClient) //verbose output

	go r.fayeClient.Listen()

	for {
		msg := <-r.channel
		data := msg.Data()
		content, _ := data["subject"]
		contentType := data["type"].(string)

		switch contentType {
		case "line.create":
			b, _ := json.Marshal(content)

			out := Message{}
			json.Unmarshal(b, &out)
			//fmt.Printf("%+v\n", out) //TODO

			r.MessageChannel <- out
			break
		case "like.create":
			b, _ := json.Marshal(content.(map[string]interface{})["line"])

			out := Message{}
			//log.Println(string(b))
			err := json.Unmarshal(b, &out)
			if err != nil {
				log.Println(err)
			}
			outt := LikeEvent{Message: out}
			//fmt.Printf("Like on %+v \n", outt.Message)

			r.LikeChannel <- outt
			break
		case "ping":
			break
		default: //TODO: see if any other types are returned
			if len(contentType) == 0 || content == nil {
				break
			}
			log.Println(contentType)
			log.Fatalln(data)

		}

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
