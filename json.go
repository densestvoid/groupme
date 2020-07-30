package groupme

import (
	"encoding/json"
	"fmt"
)

// Meta is the error type returned in the GroupMe response.
// Meant for clients that can't read HTTP status codes
type Meta struct {
	Code   HTTPStatusCode `json:"code"`
	Errors []string       `json:"errors"`
}

// Error returns the code and the error list as a string.
// Satisfies the error interface
func (m Meta) Error() string {
	return fmt.Sprintf("Error Code %d: %v", m.Code, m.Errors)
}

// Group is a GroupMe group, returned in JSON API responses
type Group struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
	// Type of group (private|public)
	Type          string        `json:"type"`
	Description   string        `json:"description"`
	ImageURL      string        `json:"image_url"`
	CreatorUserID ID            `json:"creator_user_id"`
	CreatedAt     Timestamp     `json:"created_at"`
	UpdatedAt     Timestamp     `json:"updated_at"`
	Members       []*Member     `json:"members"`
	ShareURL      string        `json:"share_url"`
	Messages      GroupMessages `json:"messages"`
}

// GroupMessages is a Group field, only returned in Group JSON API responses
type GroupMessages struct {
	Count                uint           `json:"count"`
	LastMessageID        ID             `json:"last_message_id"`
	LastMessageCreatedAt Timestamp      `json:"last_message_created_at"`
	Preview              MessagePreview `json:"preview"`
}

// MessagePreview is a GroupMessages field, only returned in Group JSON API responses.
// Abbreviated form of Message type
type MessagePreview struct {
	Nickname    string        `json:"nickname"`
	Text        string        `json:"text"`
	ImageURL    string        `json:"image_url"`
	Attachments []*Attachment `json:"attachments"`
}

// GetMemberByUserID gets the group member by their UserID,
// nil if no member matches
func (g Group) GetMemberByUserID(userID ID) *Member {
	for _, member := range g.Members {
		if member.UserID == userID {
			return member
		}
	}

	return nil
}

// GetMemberByNickname gets the group member by their Nickname,
// nil if no member matches
func (g Group) GetMemberByNickname(nickname string) *Member {
	for _, member := range g.Members {
		if member.Nickname == nickname {
			return member
		}
	}

	return nil
}

func (g Group) String() string {
	return marshal(&g)
}

// Member is a GroupMe group member, returned in JSON API responses
type Member struct {
	ID           ID     `json:"id"`
	UserID       ID     `json:"user_id"`
	Nickname     string `json:"nickname"`
	Muted        bool   `json:"muted"`
	ImageURL     string `json:"image_url"`
	AutoKicked   bool   `json:"autokicked"`
	AppInstalled bool   `json:"app_installed"`
	GUID         string `json:"guid"`
}

func (m Member) String() string {
	return marshal(&m)
}

// Message is a GroupMe group message, returned in JSON API responses
type Message struct {
	ID             ID         `json:"id"`
	SourceGUID     string     `json:"source_guid"`
	CreatedAt      Timestamp  `json:"created_at"`
	GroupID        ID         `json:"group_id"`
	UserID         ID         `json:"user_id"`
	BotID          ID         `json:"bot_id"`
	SenderID       ID         `json:"sender_id"`
	SenderType     SenderType `json:"sender_type"`
	System         bool       `json:"system"`
	Name           string     `json:"name"`
	RecipientID    ID         `json:"recipient_id"`
	ConversationID ID         `json:"conversation_id"`
	AvatarURL      string     `json:"avatar_url"`
	// Maximum length of 1000 characters
	Text string `json:"text"`
	// Must be an image service URL (i.groupme.com)
	ImageURL    string        `json:"image_url"`
	FavoritedBy []string      `json:"favorited_by"`
	Attachments []*Attachment `json:"attachments"`
}

func (m Message) String() string {
	return marshal(&m)
}

type SenderType string

// SenderType constants
const (
	SenderType_User   SenderType = "user"
	SenderType_Bot    SenderType = "bot"
	SenderType_System SenderType = "system"
)

type AttachmentType string

// AttachmentType constants
const (
	Mentions AttachmentType = "mentions"
	Image    AttachmentType = "image"
	Location AttachmentType = "location"
	Emoji    AttachmentType = "emoji"
)

// Attachment is a GroupMe message attachment, returned in JSON API responses
type Attachment struct {
	Type        AttachmentType `json:"type"`
	Loci        [][]int        `json:"loci"`
	UserIDs     []ID           `json:"user_ids"`
	URL         string         `json:"url"`
	Name        string         `json:"name"`
	Latitude    string         `json:"lat"`
	Longitude   string         `json:"lng"`
	Placeholder string         `json:"placeholder"`
	Charmap     [][]int        `json:"charmap"`
}

func (a Attachment) String() string {
	return marshal(&a)
}

// User is a GroupMe user, returned in JSON API responses
type User struct {
	ID          ID          `json:"id"`
	PhoneNumber PhoneNumber `json:"phone_number"`
	ImageURL    string      `json:"image_url"`
	Name        string      `json:"name"`
	CreatedAt   Timestamp   `json:"created_at"`
	UpdatedAt   Timestamp   `json:"updated_at"`
	AvatarURL   string      `json:"avatar_url"`
	Email       string      `json:"email"`
	SMS         bool        `json:"sms"`
}

func (u User) String() string {
	return marshal(&u)
}

// Chat is a GroupMe direct message conversation between two users,
// returned in JSON API responses
type Chat struct {
	CreatedAt     Timestamp `json:"created_at"`
	UpdatedAt     Timestamp `json:"updated_at"`
	LastMessage   *Message  `json:"last_message"`
	MessagesCount int       `json:"messages_count"`
	OtherUser     User      `json:"other_user"`
}

func (c Chat) String() string {
	return marshal(&c)
}

type Bot struct {
	BotID          ID     `json:"bot_id"`
	GroupID        ID     `json:"group_id"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	CallbackURL    string `json:"callback_url"`
	DMNotification bool   `json:"dm_notification"`
}

func (b Bot) String() string {
	return marshal(&b)
}

type Block struct {
	UserID        ID        `json:"user_id"`
	BlockedUserID ID        `json:"blocked_user_id"`
	CreatedAT     Timestamp `json:"created_at"`
}

func (b Block) String() string {
	return marshal(&b)
}

// Superficially increases test coverage
func marshal(i interface{}) string {
	bytes, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return ""
	}

	return string(bytes)
}
