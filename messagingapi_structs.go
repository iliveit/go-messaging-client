//Package messagingapi implements the iliveit Messaging API
package messagingapi

import (
	
)

// MessagingAPI is the primary object to work with the API
type MessagingAPI struct {
	config APIConfig
	
} 

// APIWebRequest handles all API communication
type APIWebRequest struct {
	config APIConfig 
	Url string
	Method string
	Data string
}

// APIConfig is the API configuration passed when creating new API instance
type APIConfig struct {
	Endpoint string
    AccessToken string
}

// NewMessage is the wrapper struct to submit a new message to the API
type NewMessage struct {
	// The action to be taken by the API
	Action int
	// The Operator this message belongs to
	MVNOID int
	// The data for that is required for building or submitting the message
	Data interface{}
	// The Campaign id for tracking purposes (optional)
	Campaign string
	// A URL where SMS replies to this message should be POSTed. Applies only to SMS messages (optional)
	PostbackReplyUrl string
	// A URL where status updates for this message should be POSTed (optional)
	PostbackStatusUrl string
	
	Error string
}

// NewMessageResult is returned with all message submit requests
type NewMessageResult struct {
	 MessageID string
}

// APIResult is the result returned with any API request
type APIResult struct {
	StatusCode uint32
	StatusDescription string
	MessageResult NewMessageResult
	//ScrubResult ScrubResult;
	//MessageStatus MessageStatus;
	//ArchivedMessage ArchivedMessage;
}

// WebRequestResponse is returned from a failed webrequest
type WebRequestResponse struct {
	Error string
}

// SubmitMMSMessageData holds information for submitting an MMS
type SubmitMMSMessageData struct {
	// The message type, 'mms' for MMS messages
	MessageType string
	MSISDN []string
	Network string
	// List of slides to use for the message
	Slides []MMSSlide `json:"slides"`
	// The MMS message subject
	Subject string `json:"subject"`
}

// MMSSlide holds the slide information for an MMS message
type MMSSlide struct {
	// The duration of the slide in seconds
	Duration string `json:"duration"`
	//  The content for the slide
	Content []MMSSlideContent `json:"content"`
}

// MMSSlideContent is the content of an MMS slide
type MMSSlideContent struct {
	// The type of slide
	Type string `json:"type"`
	// The mime type of the content of the slide
	Mime string `json:"mime"`
	// The data used in the slide, base64 encoded
	Data string `json:"data"`
	// The name to be used for the slide content
	Name string `json:"name"`
}

// SubmitSMSMessageData holds information for submitting an SMS
type SubmitSMSMessageData struct {
	// The message type, 'mms' for MMS messages
	MessageType string
	MSISDN []string
	Network string
	// Message is the actual message content of the SMS
	Message string `json:"text"`
	// Extra digits to append to sender address, when allowed
	ExtraDigits string `json:"extra_digits"`
}

// SubmitEmailMessageData holds information for submitting an Email
type SubmitEmailMessageData struct {
	// The message type, 'mms' for MMS messages
	MessageType string
	Address []string `json:"address"`
	MSISDN []string
	Network string 
	Subject string `json:"subject"`
	HTML string `json:"html"`
	Text string `json:"text"`
}

// IncomingSMS is the structure returned from the 
// API for incoming SMS messages
type IncomingSMS struct {
	// The message ID assigned by the API when the SMS
	// was submitted using Create()
	MessageId string
	// The MSISDN who sent the message
    SourceMSISDN string
    // The MSISDN the message was sent to
    DestinationMSISDN string
    // The SMS message
    Message string
    // The extra digits supplied when the message was submitted
    ExtraDigits string
    // The amount of times the message was retried to you
    RetryCount uint32
}


const (
	APIResultStatusesOk = 0
	APIResultStatusesError = 1
	APIResultStatusesAuthFailed = 2
	APIResultStatusesInvalidMethod = 3
	APIResultStatusesAPIError = 4
	APIResultStatusesRateLimited = 5
)

const (
	APIActionTypesSubmitMMS = 1
	APIActionTypesSubmitSMS = 2
	APIActionTypesSubmitEmail = 3
	APIActionTypesArchive = 4
)

const (
	MMSContentTypeText = "text"
	MMSContentTypeImage = "image"
	MMSContentTypeVideo = "video"
	MMSContentTypeAudio = "audio"
	
)