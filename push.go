package exponent

import (
	"errors"
	"net/http"
	"strings"
)

type Priority string
type Data map[string]string
type ErrorMsg string
type Token string

const (
	// NormalPriority is a priority used in PushMessage
	NormalPriority Priority = "normal"
	// HighPriority is a priority used in PushMessage
	HighPriority Priority = "high"
	// DefaultPriority is the standard priority used in PushMessage
	DefaultPriority Priority = "default"

	// ErrorMsgDeviceNotRegistered indicates the token is invalid
	ErrorMsgDeviceNotRegistered ErrorMsg = "DeviceNotRegistered"
	// ErrorMsgTooBig indicates the message went over payload size of 4096 bytes
	ErrorMsgTooBig ErrorMsg = "MessageTooBig"
	// ErrorMsgRateExceeded indicates messages have been sent too frequently
	ErrorMsgRateExceeded ErrorMsg = "MessageRateExceeded"
	// ErrMsgMalformedToken is returned if a token does not start with 'ExponentPushToken'
	ErrMsgMalformedToken ErrorMsg = "token should start with ExponentPushToken"
)

// ParseToken returns a token and may return an error if the input token is invalid
func ParseToken(token string) (*Token, error) {
	if !strings.HasPrefix(token, "ExponentPushToken") {
		return nil, errors.New(string(ErrMsgMalformedToken))
	}
	tkn := Token(token)
	return &tkn, nil
}

func MustParseToken(token string) *Token {
	tkn, _ := ParseToken(token)
	return tkn
}

func IsPushTokenValid(token string) bool {
	_, err := ParseToken(token)
	return err == nil
}

// is an object that describes a push notification request.
type Message struct {
	// An Expo push token or an array of Expo push tokens specifying the recipient(s) of this message.
	To []*Token `json:"to"`
	// The title to display in the notification. On iOS, this is displayed only on Apple Watch.
	Title string `json:"title,omitempty"`
	// The message to display in the notification.
	Body string `json:"body"`
	// A dict of extra data to pass inside of the push notification. The total notification payload must be at most 4096 bytes.
	Data Data `json:"data,omitempty"`
	// A sound to play when the recipient receives this notification.
	// Specify "default" to play the device's default notification sound, or omit this field to play no sound.
	Sound string `json:"sound,omitempty"`
	// The number of seconds for which the message may be kept around for redelivery if it hasn't been delivered yet. Defaults to 0.
	TTL int `json:"ttl,omitempty"`
	// UNIX timestamp for when this message expires. It has the same effect as ttl, and is just an absolute timestamp instead of a relative one.
	Expiration int64 `json:"expiration,omitempty"`
	// Delivery priority of the message. Use the *Priority constants specified above.
	Priority Priority `json:"priority,omitempty"`
	// An integer representing the unread notification count.
	// This currently only affects iOS. Specify 0 to clear the badge count.
	Badge int `json:"badge,omitempty"`
	// ID of the Notification Channel through which to display this notification on Android devices.
	ChannelID string `json:"channelId,omitempty"`
}

// Response is the HTTP response returned from an Expo publish HTTP request
type Response struct {
	Data   []*MessageResponse `json:"data"`
	Errors []Data             `json:"errors"`
}

// PushResponse is a wrapper class for a push notification response.
// A successful single push notification:
//
//	{'status': 'ok'}
//
// An invalid push token
//
//	{'status': 'error',
//	 'message': '"adsf" is not a registered push notification recipient'}
type MessageResponse struct {
	MessageItem *Message
	ID          string `json:"id"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	Details     Data   `json:"details"`
}

func (r *MessageResponse) IsOk() bool {
	return r.Status == "ok"
}

// ServerError is raised when the push token server is not behaving as expected
// For example, invalid push notification arguments result in a different
// style of error. Instead of a "data" array containing errors per
// notification, an "error" array is returned.
// {"errors": [
//
//	{"code": "API_ERROR",
//	 "message": "child \"to\" fails because [\"to\" must be a string]. \"value\" must be an array."
//	}
//
// ]}
type ServerError struct {
	Message      string
	Response     *http.Response
	ResponseData *Response
	Errors       []Data
}

// NewServerError creates a new PushServerError object
func NewServerError(message string, response *http.Response,
	responseData *Response,
	errors []Data) *ServerError {
	return &ServerError{
		Message:      message,
		Response:     response,
		ResponseData: responseData,
		Errors:       errors,
	}
}

func (e *ServerError) Error() string {
	return e.Message
}
