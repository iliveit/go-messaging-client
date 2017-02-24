//Package messagingapi implements the iliveit Messaging API
package messagingapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

// New creates a new API instance using the given config
func New(config APIConfig) (*MessagingAPI, error) {

	if config.AccessToken == "" {
		return nil, errors.New("Access Token can not be blank")

	}
	if config.Endpoint == "" {
		return nil, errors.New("Endpoint can not be blank")
	}

	// Add trailing slash if not present
	if config.Endpoint[len(config.Endpoint)-1:] != "/" {
		config.Endpoint = config.Endpoint + "/"
	}

	api := MessagingAPI{
		config: config,
	}

	return &api, nil
}

// HandleErrorResponse checks status codes and completed the required fields
func HandleErrorResponse(result APIResult, statusCode int, err error) APIResult {
	if statusCode != http.StatusBadGateway {
		if statusCode == 429 {
			result.StatusCode = APIResultStatusesRateLimited
		} else if statusCode == http.StatusUnauthorized {
			result.StatusCode = APIResultStatusesAuthFailed
		} else if statusCode == http.StatusBadRequest {
			result.StatusCode = APIResultStatusesError
		} else {
			result.StatusCode = APIResultStatusesAuthFailed
		}
		result.StatusDescription = err.Error()
	} else {
		result.StatusCode = APIResultStatusesError
		result.StatusDescription = err.Error()
	}
	return result
}

// Ping makes a simple GET call to the API to check authentication
func (api *MessagingAPI) Ping() (APIResult, error) {
	result := APIResult{}
	r, err := NewAPIWebRequest(api.config, "ping", "GET", "")
	if err != nil {
		return result, err
	}

	_, statusCode, err := r.Execute()
	if err != nil {
		result = HandleErrorResponse(result, statusCode, err)
	} else {
		result.StatusCode = APIResultStatusesOk
		result.StatusDescription = "Ok"
	}

	return result, nil
}

// Create requests a new message to be submitted via the API
func (api *MessagingAPI) Create(message NewMessage) (APIResult, error) {
	result := APIResult{}
	err := message.Validate()
	if err != nil {
		return result, err
	} else {

		jsonBytes, err := json.Marshal(message)
		if err != nil {
			return result, err
		}

		r, err := NewAPIWebRequest(api.config, "message/send", "POST", string(jsonBytes))
		if err != nil {
			return result, err
		}

		responseBody, statusCode, err := r.Execute()
		if err != nil {
			result = HandleErrorResponse(result, statusCode, err)
		} else {

			var newMessageResult NewMessageResult
			err = json.Unmarshal([]byte(responseBody), &newMessageResult)
			if err != nil {
				result.StatusCode = APIResultStatusesError
				result.StatusDescription = "Unable to unmarshal result from API"
			}
			result.MessageResult = newMessageResult
			result.StatusCode = APIResultStatusesOk
			result.StatusDescription = "Ok"
		}

		return result, nil
	}

	return result, nil
}

// Resend resubmits a message
func (api *MessagingAPI) Resend(resendRequest ResendMessageRequest) (APIResult, error) {
	result := APIResult{}
	jsonBytes, err := json.Marshal(resendRequest)
	if err != nil {
		return result, err
	}

	r, err := NewAPIWebRequest(api.config, "message/resend", "POST", string(jsonBytes))
	if err != nil {
		return result, err
	}
	responseBody, statusCode, err := r.Execute()
	if err != nil {
		result = HandleErrorResponse(result, statusCode, err)
	} else {

		var newMessageResult NewMessageResult
		err = json.Unmarshal([]byte(responseBody), &newMessageResult)
		if err != nil {
			result.StatusCode = APIResultStatusesError
			result.StatusDescription = "Unable to unmarshal result from API"
		}
		result.MessageResult = newMessageResult
		result.StatusCode = APIResultStatusesOk
		result.StatusDescription = "Ok"
	}

	return result, nil
}

// Create requests a new approval request to be submitted via the API
func (api *MessagingAPI) CreateApproval(approvalRequest ApprovalRequest) (APIResult, error) {
	result := APIResult{}
	jsonBytes, err := json.Marshal(approvalRequest)
	if err != nil {
		return result, err
	}

	r, err := NewAPIWebRequest(api.config, "approval/create", "POST", string(jsonBytes))
	if err != nil {
		return result, err
	}

	responseBody, statusCode, err := r.Execute()
	if err != nil {
		result = HandleErrorResponse(result, statusCode, err)
	} else {

		var requestResult ApprovalRequestResult
		err = json.Unmarshal([]byte(responseBody), &requestResult)
		if err != nil {
			result.StatusCode = APIResultStatusesError
			result.StatusDescription = "Unable to unmarshal result from API"
		}
		result.RequestResult = requestResult
		result.StatusCode = APIResultStatusesOk
		result.StatusDescription = "Ok"
	}
	return result, nil
}

// UpdateApproval requests an approval should be updated
func (api *MessagingAPI) UpdateApproval(updateRequest ApprovalUpdateRequest) (APIResult, error) {
	result := APIResult{}
	jsonBytes, err := json.Marshal(updateRequest)
	if err != nil {
		return result, err
	}

	r, err := NewAPIWebRequest(api.config, "approval/update", "PUT", string(jsonBytes))
	if err != nil {
		return result, err
	}

	responseBody, statusCode, err := r.Execute()
	if err != nil {
		result = HandleErrorResponse(result, statusCode, err)
	} else {

		var requestResult ApprovalRequestResult
		err = json.Unmarshal([]byte(responseBody), &requestResult)
		if err != nil {
			result.StatusCode = APIResultStatusesError
			result.StatusDescription = "Unable to unmarshal result from API"
		}
		result.RequestResult = requestResult
		result.StatusCode = APIResultStatusesOk
		result.StatusDescription = "Ok"
	}
	return result, nil
}

// GetMessageStatus makes a simple GET call to the API to get the status of a message
func (api *MessagingAPI) GetMessageStatus(message_id string) (APIResult, error) {
	result := APIResult{}

	if message_id == "" {
		return result, errors.New("Message ID must not be blank")
	}

	r, err := NewAPIWebRequest(api.config, "message/"+message_id+"/status", "GET", "")
	if err != nil {
		return result, err
	}

	responseBody, statusCode, err := r.Execute()
	if err != nil {
		result = HandleErrorResponse(result, statusCode, err)
	} else {

		var status StatusResult
		err = json.Unmarshal([]byte(responseBody), &status)
		if err != nil {
			result.StatusCode = APIResultStatusesError
			result.StatusDescription = "Unable to unmarshal result from API"
		}
		result.MessageStatus = status
		result.StatusCode = APIResultStatusesOk
		result.StatusDescription = "Ok"
	}
	return result, nil
}

// Generate requests a new video to be generated via the API and AfterAction to
// be executed
func (api *MessagingAPI) Generate(request BuildRequest) (APIResult, error) {
	result := APIResult{}
	err := request.Validate()
	if err != nil {
		return result, err
	} else {
		messageJson, err := request.Package()
		if err != nil {
			return result, err
		}

		r, err := NewAPIWebRequest(api.config, "generate/video", "POST", messageJson)
		if err != nil {
			return result, err
		}

		responseBody, statusCode, err := r.Execute()
		if err != nil {
			result = HandleErrorResponse(result, statusCode, err)
		} else {

			var newMessageResult NewMessageResult
			err = json.Unmarshal([]byte(responseBody), &newMessageResult)
			if err != nil {
				result.StatusCode = APIResultStatusesError
				result.StatusDescription = "Unable to unmarshal result from API"
			}
			result.MessageResult = newMessageResult
			result.StatusCode = APIResultStatusesOk
			result.StatusDescription = "Ok"
		}

		return result, nil
	}

	return result, nil
}

// GetMSISDNScrub retrieves the MSISDN's handset information
func (api *MessagingAPI) GetMSISDNScrub(msisdn string) (APIResult, error) {
	result := APIResult{}

	if msisdn == "" {
		return result, errors.New("msisdn must not be blank")
	}

	r, err := NewAPIWebRequest(api.config, "scrub/"+msisdn, "GET", "")
	if err != nil {
		return result, err
	}

	responseBody, statusCode, err := r.Execute()
	if err != nil {
		result = HandleErrorResponse(result, statusCode, err)
	} else {

		var scrubres ScrubResult
		err = json.Unmarshal([]byte(responseBody), &scrubres)
		if err != nil {
			result.StatusCode = APIResultStatusesError
			result.StatusDescription = "Unable to unmarshal result from API"
		}

		result.ScrubResult = scrubres
		result.StatusCode = APIResultStatusesOk
		result.StatusDescription = "Ok"
	}
	return result, nil
}
