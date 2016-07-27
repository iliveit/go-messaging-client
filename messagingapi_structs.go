//Package messagingapi implements the iliveit Messaging API
package messagingapi

import (
	"time"
)

// MessagingAPI is the primary object to work with the API
type MessagingAPI struct {
	config APIConfig
}

// APIWebRequest handles all API communication
type APIWebRequest struct {
	config APIConfig
	Url    string
	Method string
	Data   string
}

// APIConfig is the API configuration passed when creating new API instance
type APIConfig struct {
	Endpoint    string
	AccessToken string
}

// NewMessage is the wrapper struct to submit a new message to the API
type NewMessage struct {
	// The action to be taken by the API
	Action int
	// The Operator this message belongs to
	MVNOID int
	// If part of an approval batch, the ID
	ApprovalBatch uint32
	// The data for that is required for building or submitting the message
	Data interface{}
	// The Campaign id for tracking purposes (optional)
	Campaign string
	// A URL where SMS replies to this message should be POSTed. Applies only to SMS messages (optional)
	PostbackReplyUrl string
	// A URL where status updates for this message should be POSTed (optional)
	PostbackStatusUrl string
	// What updates you want to receive in the postback
	// Possible values of "build", "submit", "archive", "sent", "delivery"
	// comma delimited - i.e. "build,submit,delivery"
	PostbackStatusTypes string

	Error string
}

// NewMessageResult is returned with all message submit requests
type NewMessageResult struct {
	MessageID string
}

// APIResult is the result returned with any API request
type APIResult struct {
	StatusCode        uint32
	StatusDescription string
	MessageResult     NewMessageResult
	MessageStatus     StatusResult
	RequestResult     ApprovalRequestResult
	ScrubResult       ScrubResult
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
	MSISDN      []string
	Network     string
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
	MSISDN      []string
	Network     string
	// Message is the actual message content of the SMS
	Message string `json:"text"`
	// Extra digits to append to sender address, when allowed
	ExtraDigits string `json:"extra_digits"`
}

// SubmitEmailMessageData holds information for submitting an Email
type SubmitEmailMessageData struct {
	// The message type, 'mms' for MMS messages
	MessageType string
	Address     []string `json:"address"`
	MSISDN      []string
	Network     string
	Subject     string            `json:"subject"`
	HTML        string            `json:"html"`
	Text        string            `json:"text"`
	Attachments []EmailAttachment `json:"attachments"`
}

// EmailAttachment holds the structure of the attachment data
type EmailAttachment struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
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

// StatusResult is returned from the API when message/:id/status
// is called
type StatusResult struct {
	Type                       string    `json:"type"`
	Campaign                   string    `json:"campaign,omitempty"`
	MessageID                  string    `json:"message_id"`
	Template                   uint32    `json:"template"`
	Network                    string    `json:"network"`
	MSISDN                     string    `json:"msisdn"`
	Email                      string    `json:"email"`
	MVNO                       uint32    `json:"mvno"`
	DateReceived               time.Time `json:"date_received"`
	BuildStatus                uint32    `json:"build_status,omitempty"`
	BuildStatusDescription     string    `json:"build_status_description,omitempty"`
	BuildTimestamp             time.Time `json:"build_timestamp,omitempty"`
	ArchiveStatus              uint32    `json:"archive_status,omitempty"`
	ArchiveStatusDescription   string    `json:"archive_status_description,omitempty"`
	ArchiveTimestamp           time.Time `json:"archive_timestamp,omitempty"`
	SubmitStatus               uint32    `json:"submit_status,omitempty"`
	SubmitStatusDescription    string    `json:"submit_status_description,omitempty"`
	SubmitTimestamp            time.Time `json:"submit_timestamp,omitempty"`
	SentStatus                 uint32    `json:"sent_status,omitempty"`
	SentStatusDescription      string    `json:"sent_status_description,omitempty"`
	SentTimestamp              time.Time `json:"sent_timestamp,omitempty"`
	DeliveredStatus            uint32    `json:"delivered_status,omitempty"`
	DeliveredStatusDescription string    `json:"delivered_status_description,omitempty"`
	DeliveredTimestamp         time.Time `json:"delivered_timestamp,omitempty"`
	PostbackType               string    `json:"postback_type,omitempty"`
}

type ScrubResult struct {
	Network      string        `json:"network"`
	MSISDN       string        `json:"msisdn"`
	HandsetMake  string        `json:"handset_make"`
	HandsetModel string        `json:"handset_model"`
	AllowSend    string        `json:allow_send`
	ScreenSize   ScreenSizeObj `json:"screen_size"`
	ErrorCode    string        `json:"error_code"`
	Error        string        `json:"error"`
}

type ScreenSizeObj struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

// The person appprovals should be sent to
type ApprovalPerson struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	MSISDN string `json:"msisdn"`
	// The user's unique hash, leave blank for autogenerate
	Hash string `json:"hash,omitempty"`
}

// ApprovalRequest is the scrusture for creating approval batches
type ApprovalRequest struct {
	// The action type of the messages,
	// either ActionSubmitMMS, ActionSubmitSMS or ActionSubmitEmail
	ActionType   uint32
	MVNOID       uint32
	Name         string
	MaxApprovals uint32
	// Internal users to send approval messages to
	InternalPeople []ApprovalPerson
	// Clients to send approval messages to
	ExternalPeople []ApprovalPerson
	// The link to send with approval messages, leave blank to
	// autogenerate
	Link string `json:",omitempty"`
	// If this is linked to another approval, the approvals
	// will not be sent until all linked has been marked ready
	LinkedApproval uint32 `json:"linked_approval,omitempty"`
}

// ApprovalRequestResult The result of the approval request
type ApprovalRequestResult struct {
	// The batch that was created
	BatchID uint32
}

type ApprovalUpdateRequest struct {
	// The batch ID to update
	BatchID uint32
	// State should be one of ApprovalRequestState*
	State uint32
	// CSV data for a report, blank if no needed
	Report []string
}

const (
	ApprovalBatchStateWaitingData  = 1
	ApprovalBatchStateDataReceived = 2
	ApprovalBatchStateApprovalSent = 3
	ApprovalBatchStateApproved     = 4
	ApprovalBatchStateDeclined     = 5
	ApprovalBatchStateSent         = 6
)

const (
	APIResultStatusesOk            = 0
	APIResultStatusesError         = 1
	APIResultStatusesAuthFailed    = 2
	APIResultStatusesInvalidMethod = 3
	APIResultStatusesAPIError      = 4
	APIResultStatusesRateLimited   = 5
)

const (
	APIActionTypesSubmitMMS    = 1
	APIActionTypesSubmitSMS    = 2
	APIActionTypesSubmitEmail  = 3
	APIActionTypesArchive      = 4
	APIActionTypesArchiveMMS   = 8
	APIActionTypesArchiveSMS   = 9
	APIActionTypesArchiveEmail = 10
)

const (
	MMSContentTypeText  = "text"
	MMSContentTypeImage = "image"
	MMSContentTypeVideo = "video"
	MMSContentTypeAudio = "audio"
)

// Message Status Types
const (
	// Received by scheduler
	MessageStatusReceived = 1 + iota
	// Submitted to Beam
	MessageStatusSubmitted
	// Submission to Beam failed
	MessageStatusSubmitFailed
	// Archived
	MessageArchived
	// Archived because of network action
	MessageArchivedNetworkRule
	// Invalid JSON output
	MessageStatusFailedMarshal
	// Failed due to bad JSON
	MessageStatusFailedUnmarshal
	// When the scrub service is unavailable
	MessageStatusScrubUnavailable
	// Status when the MSISDN is not MMS capable
	MessageStatusNotMMS
	// When failed to build message from a template
	MessageStatusTemplateError
	// No network found for this message
	MessageStatusNoNetwork
	// Error submitting to renderfarm
	MessageStatusRenderFarmFailed
	// Sent to for rendering
	MessageStatusRenderFarmSent
	// Success render
	MessageStatusRenderFarmSuccess
	// Unable to read file
	MessageStatusReadFileFailed
	// Invalid data received
	MessageInvalid
	// No handset information available for MSISDN
	MessageStatusNoHandsetInfo
	// MSISDN is not provisioned for MMS
	MessageStatusNotMMSProvisioned
)
