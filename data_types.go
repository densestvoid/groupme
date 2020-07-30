package groupme

import (
	"regexp"
	"time"
)

// GroupMe documentation: https://dev.groupme.com/docs/responses

// ID is an unordered alphanumeric string
type ID string

// Treated as a constant
var alphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// Valid checks if the ID string is alpha numeric
func (id ID) Valid() bool {
	return alphaNumericRegex.MatchString(string(id))
}

func (id ID) String() string {
	return string(id)
}

// Timestamp is the number of seconds since the UNIX epoch
type Timestamp uint64

// FromTime returns the time.Time as a Timestamp
func FromTime(t time.Time) Timestamp {
	return Timestamp(t.Unix())
}

// ToTime returns the Timestamp as a UTC Time
func (t Timestamp) ToTime() time.Time {
	return time.Unix(int64(t), 0).UTC()
}

// String returns the Timestamp in the default time.Time string format
func (t Timestamp) String() string {
	return t.ToTime().String()
}

// PhoneNumber is the country code plus the number of the user
type PhoneNumber string

// Treated as a constant
var phoneNumberRegex = regexp.MustCompile(`^\+[0-9]+ [0-9]{10}$`)

// Valid checks if the ID string is alpha numeric
func (pn PhoneNumber) Valid() bool {
	return phoneNumberRegex.MatchString(string(pn))
}

func (pn PhoneNumber) String() string {
	return string(pn)
}

// StatusCodes are returned by HTTP requests in
// the header and the json "meta" field
type HTTPStatusCode int

// Text used as constant name
const (
	HTTP_Ok                  HTTPStatusCode = 200
	HTTP_Created             HTTPStatusCode = 201
	HTTP_NoContent           HTTPStatusCode = 204
	HTTP_NotModified         HTTPStatusCode = 304
	HTTP_BadRequest          HTTPStatusCode = 400
	HTTP_Unauthorized        HTTPStatusCode = 401
	HTTP_Forbidden           HTTPStatusCode = 403
	HTTP_NotFound            HTTPStatusCode = 404
	HTTP_EnhanceYourCalm     HTTPStatusCode = 420
	HTTP_InternalServerError HTTPStatusCode = 500
	HTTP_BadGateway          HTTPStatusCode = 502
	HTTP_ServiceUnavailable  HTTPStatusCode = 503
)

// String returns the description of the status code according to GroupMe
func (c HTTPStatusCode) String() string {
	return map[HTTPStatusCode]string{
		HTTP_Ok:                  "success",
		HTTP_Created:             "resource was created successfully",
		HTTP_NoContent:           "resource was deleted successfully",
		HTTP_NotModified:         "no new data to return",
		HTTP_BadRequest:          "invalid format or data specified in the request",
		HTTP_Unauthorized:        "authentication credentials missing or incorrect",
		HTTP_Forbidden:           "request refused due to update limits",
		HTTP_NotFound:            "URI is invalid or resource does not exist",
		HTTP_EnhanceYourCalm:     "application is being rate limited",
		HTTP_InternalServerError: "something unexpected occurred",
		HTTP_BadGateway:          "GroupMe is down or being upgraded",
		HTTP_ServiceUnavailable:  "servers are overloaded, try again later",
	}[c]
}
