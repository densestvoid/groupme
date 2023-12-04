package groupme

import (
	"net/http"
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

// StatusText returns a text for the HTTP status code (according to GroupMe). It returns the empty string if the code is unknown.
func HTTPStatusText(code int) string {
	return map[int]string{
		http.StatusOK:                  "Success!",                                                                                                                                                                   // OK
		http.StatusCreated:             "Resource was created successfully.",                                                                                                                                         // Created
		http.StatusNoContent:           "Resource was deleted successfully.",                                                                                                                                         // No Content
		http.StatusNotModified:         "There was no new data to return.",                                                                                                                                           // Not Modified
		http.StatusBadRequest:          "Returned when an invalid format or invalid data is specified in the request.",                                                                                               // Bad Request
		http.StatusUnauthorized:        "Authentication credentials were missing or incorrect.",                                                                                                                      // Unauthorized
		http.StatusForbidden:           "The request is understood, but it has been refused. An accompanying error message will explain why. This code is used when requests are being denied due to update limits.", // Forbidden
		http.StatusNotFound:            "The URI requested is invalid or the resource requested, such as a user, does not exists.",                                                                                   // Not Found
		420:                            "Returned when you are being rate limited. Chill the heck out.",                                                                                                              // Enhance Your Calm
		http.StatusInternalServerError: "Something unexpected occurred. GroupMe will be notified.",                                                                                                                   // Internal Server Error
		http.StatusBadGateway:          "GroupMe is down or being upgraded.",                                                                                                                                         // Bad Gateway
		http.StatusServiceUnavailable:  "The GroupMe servers are up, but overloaded with requests. Try again later.",                                                                                                 // Service Unavailable
	}[code]
}
