// Package groupme defines a client capable of executing API commands for the GroupMe chat service
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
var phoneNumberRegex = regexp.MustCompile(`^\+\d+ \d{10}$`)

// Valid checks if the ID string is alpha numeric
func (pn PhoneNumber) Valid() bool {
	return phoneNumberRegex.MatchString(string(pn))
}

func (pn PhoneNumber) String() string {
	return string(pn)
}

// HTTPStatusCode are returned by HTTP requests in
// the header and the json "meta" field
type HTTPStatusCode int

// Text used as constant name
const (
	HTTPOk                  HTTPStatusCode = 200
	HTTPCreated             HTTPStatusCode = 201
	HTTPNoContent           HTTPStatusCode = 204
	HTTPNotModified         HTTPStatusCode = 304
	HTTPBadRequest          HTTPStatusCode = 400
	HTTPUnauthorized        HTTPStatusCode = 401
	HTTPForbidden           HTTPStatusCode = 403
	HTTPNotFound            HTTPStatusCode = 404
	HTTPEnhanceYourCalm     HTTPStatusCode = 420
	HTTPInternalServerError HTTPStatusCode = 500
	HTTPBadGateway          HTTPStatusCode = 502
	HTTPServiceUnavailable  HTTPStatusCode = 503
)

// String returns the description of the status code according to GroupMe
func (c HTTPStatusCode) String() string {
	return map[HTTPStatusCode]string{
		HTTPOk:                  "success",
		HTTPCreated:             "resource was created successfully",
		HTTPNoContent:           "resource was deleted successfully",
		HTTPNotModified:         "no new data to return",
		HTTPBadRequest:          "invalid format or data specified in the request",
		HTTPUnauthorized:        "authentication credentials missing or incorrect",
		HTTPForbidden:           "request refused due to update limits",
		HTTPNotFound:            "URI is invalid or resource does not exist",
		HTTPEnhanceYourCalm:     "application is being rate limited",
		HTTPInternalServerError: "something unexpected occurred",
		HTTPBadGateway:          "GroupMe is down or being upgraded",
		HTTPServiceUnavailable:  "servers are overloaded, try again later",
	}[c]
}
