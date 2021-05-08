package groupme

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func init() {

	RealTimeHandlers = make(map[string]func(r *PushSubscription, channel string, data ...interface{}))

	//Base Handlers on user channel
	RealTimeHandlers["direct_message.create"] = func(r *PushSubscription, channel string, data ...interface{}) {
		b, _ := json.Marshal(data[0])
		out := Message{}
		_ = json.Unmarshal(b, &out)

		//maybe something with API versioning
		out.ConversationID = out.ChatID

		if out.UserID.String() == "system" {
			event := struct {
				Event struct {
					Kind string `json:"type"`
					Data interface{}
				}
			}{}

			err := json.Unmarshal(b, &event)
			if err != nil {
				fmt.Println(err)
			}
			rawData, _ := json.Marshal(event.Event.Data)
			handler, ok := RealTimeSystemHandlers[event.Event.Kind]
			if !ok {
				log.Println("Unable to handle system message of type", event.Event.Kind)
				return
			}

			id := out.GroupID
			if len(id) == 0 {
				id = out.ConversationID
			}

			handler(r, channel, id, rawData)
			return
		}

		for _, h := range r.handlers {
			if h, ok := h.(HandlerText); ok {
				h.HandleTextMessage(out)
			}
		}
	}

	RealTimeHandlers["line.create"] = RealTimeHandlers["direct_message.create"]

	RealTimeHandlers["like.create"] = func(r *PushSubscription, channel string, data ...interface{}) { //should be an associated chatEvent
	}

	RealTimeHandlers["membership.create"] = func(r *PushSubscription, channel string, data ...interface{}) {
		c, _ := data[0].(map[string]interface{})
		id, _ := c["id"].(string)

		for _, h := range r.handlers {
			if h, ok := h.(HandlerMembership); ok {
				h.HandleJoin(ID(id))
			}
		}

	}

	//following are for each chat
	RealTimeHandlers["favorite"] = func(r *PushSubscription, channel string, data ...interface{}) {
		c, ok := data[0].(map[string]interface{})
		if !ok {
			fmt.Println(data, "err")
			return
		}
		e, ok := c["line"]
		if !ok {
			fmt.Println(data, "err")
			return
		}
		d, _ := json.Marshal(e)
		msg := Message{}
		_ = json.Unmarshal(d, &msg)
		for _, h := range r.handlers {
			if h, ok := h.(HandlerLike); ok {
				h.HandleLike(msg)
			}
		}
	}

	//following are for messages from system (administrative/settings changes)
	RealTimeSystemHandlers = make(map[string]func(r *PushSubscription, channel string, id ID, rawData []byte))

	RealTimeSystemHandlers["membership.nickname_changed"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {
		thing := struct {
			Name string
			User struct {
				ID int
			}
		}{}
		_ = json.Unmarshal(rawData, &thing)

		for _, h := range r.handlers {
			if h, ok := h.(HandleMemberNewNickname); ok {
				h.HandleNewNickname(id, ID(strconv.Itoa(thing.User.ID)), thing.Name)
			}
		}

	}

	RealTimeSystemHandlers["membership.avatar_changed"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {
		content := struct {
			AvatarURL string `json:"avatar_url"`
			User      struct {
				ID int
			}
		}{}
		_ = json.Unmarshal(rawData, &content)

		for _, h := range r.handlers {
			if h, ok := h.(HandleMemberNewAvatar); ok {
				h.HandleNewAvatarInGroup(id, ID(strconv.Itoa(content.User.ID)), content.AvatarURL)
			}
		}

	}

	RealTimeSystemHandlers["membership.announce.added"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {
		data := struct {
			Added []Member `json:"added_users"`
		}{}
		_ = json.Unmarshal(rawData, &data)
		for _, h := range r.handlers {
			if h, ok := h.(HandleMembers); ok {
				h.HandleMembers(id, data.Added, true)
			}
		}
	}

	RealTimeSystemHandlers["membership.notifications.removed"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {
		data := struct {
			Added Member `json:"removed_user"`
		}{}
		_ = json.Unmarshal(rawData, &data)
		for _, h := range r.handlers {
			if h, ok := h.(HandleMembers); ok {
				h.HandleMembers(id, []Member{data.Added}, false)
			}
		}

	}

	RealTimeSystemHandlers["membership.name_change"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {

		data := struct {
			Name string
		}{}
		_ = json.Unmarshal(rawData, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupName); ok {
				h.HandleGroupName(id, data.Name)
			}
		}
	}

	RealTimeSystemHandlers["group.name_change"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {

		data := struct {
			Name string
		}{}
		_ = json.Unmarshal(rawData, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupName); ok {
				h.HandleGroupName(id, data.Name)
			}
		}
	}

	RealTimeSystemHandlers["group.topic_change"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {

		data := struct {
			Topic string
		}{}
		_ = json.Unmarshal(rawData, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupTopic); ok {
				h.HandleGroupTopic(id, data.Topic)
			}
		}
	}

	RealTimeSystemHandlers["group.avatar_change"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {
		data := struct {
			AvatarURL string `json:"avatar_url"`
		}{}
		_ = json.Unmarshal(rawData, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupAvatar); ok {
				h.HandleGroupAvatar(id, data.AvatarURL)
			}
		}
	}

	RealTimeSystemHandlers["group.like_icon_set"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {
		data := struct {
			LikeIcon struct {
				PackID    int `json:"pack_id"`
				PackIndex int `json:"pack_index"`
				Type      string
			} `json:"like_icon"`
		}{}
		_ = json.Unmarshal(rawData, &data)

		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupLikeIcon); ok {
				h.HandleLikeIcon(id, data.LikeIcon.PackID, data.LikeIcon.PackIndex, data.LikeIcon.Type)
			}
		}
	}

	RealTimeSystemHandlers["group.like_icon_removed"] = func(r *PushSubscription, channel string, id ID, rawData []byte) {
		for _, h := range r.handlers {
			if h, ok := h.(HandleGroupLikeIcon); ok {
				h.HandleLikeIcon(id, 0, 0, "")
			}
		}
	}

}
