//Package messagingapi implements the iliveit Messaging API
package messagingapi

import (
	"encoding/json"
	"errors"
)

type BuildRequest struct {

	// MVNOID this message belongs to
	MVNOID int
	//Data is the object that is required for building the message
	Data interface{}
	// Campaign is the Campaign id for tracking purposes (optional)
	Campaign string
	// If part of an approval batch, the ID
	ApprovalBatch uint32
	// BuildTemplate is the tempate to be used when building the message
	BuildTemplate int
	// The action to be taken by the API after the message has been built
	// Listed as APIActionTypes* constants
	AfterBuildAction int
	// The data required for the AfterBuildAction
	AfterBuildData interface{}
	// A URL where status updates for this message should be POSTed (optional)
	PostbackStatusUrl string
	// Allows you to force a size of a build. When ForcedSize is 'Both', AfterBuildAction must be archive
	ForcedSize string
	// PostBackStatusTypes determines which updates you will receive in the postback
	// Possible values of "build", "submit", "archive", "sent", "delivery"
	// comma delimited - i.e. "build,submit,delivery"
	PostBackStatusTypes string
	// Error is the last error that occurred within Validate()
	Error string
}

// Package marshals the BuildRequest and returns the JSON string and/or error
func (this *BuildRequest) Package() (string, error) {

	// If this.Data is a plain string, don't marshal
	if _, ok := this.Data.(string); ok == false {
		jsonBytes, err := json.Marshal(this.Data)
		if err != nil {
			return "", err
		}
		this.Data = string(jsonBytes)
	}

	jsonBytes, err := json.Marshal(this.AfterBuildData)
	if err != nil {
		return "", err
	}
	this.AfterBuildData = string(jsonBytes)

	jsonBytes, err = json.Marshal(this)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// Validate validates the message before sending via the API
// returns true if valid, false otherwise
func (this *BuildRequest) Validate() error {
	var err error
	if this.AfterBuildAction == 0 {
		err = errors.New("AfterBuildAction must be set")
	}
	if this.MVNOID == 0 {
		err = errors.New("MVNOID must be set and not zero")
	}
	if this.Data == nil {
		err = errors.New("A build request must have data")
	}
	if this.AfterBuildAction != APIActionTypesArchive {

		if this.AfterBuildData == nil {
			err = errors.New("If the AfterBuildAction is not Archive, AfterBuildData needs to be specified")
		}

		if this.AfterBuildAction == APIActionTypesSubmitMMS {
			if _, ok := this.AfterBuildData.(SubmitMMSMessageData); ok == false {
				err = errors.New("Using AfterAction of SubmitMMS requires AfterBuildData to be of type SubmitMMSMessageData")
			}
		}

	}
	if this.BuildTemplate == 0 {
		err = errors.New("A build template must be selected")
	}
	return err
}
