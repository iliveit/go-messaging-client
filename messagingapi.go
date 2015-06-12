//Package messagingapi implements the iliveit Messaging API
package messagingapi

import (
	"errors"
	"net/http"
	"encoding/json"
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
    if config.Endpoint[len(config.Endpoint) - 1:] != "/" {
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
	r, err :=  NewAPIWebRequest(api.config, "ping", "GET", "")
	if err != nil {
		return result, err
	}
	
	_, statusCode, err := r.Execute()
	if err != nil {
		result = HandleErrorResponse(result, statusCode, err)
	} else {
		result.StatusCode = APIResultStatusesOk;
		result.StatusDescription = "Ok";
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
		
		r, err :=  NewAPIWebRequest(api.config, "message/send", "POST", string(jsonBytes))
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
				result.StatusCode = APIResultStatusesError;
				result.StatusDescription = "Unable to unmarshal result from API";
			}
			result.MessageResult = newMessageResult
			result.StatusCode = APIResultStatusesOk;
			result.StatusDescription = "Ok";
		}
		
		return result, nil
	}
	
	return result, nil
}

// GetMessageStatus makes a simple GET call to the API to get the status of a message
func (api *MessagingAPI) GetMessageStatus(message_id string) (APIResult, error) {
	result := APIResult{}
	
	if message_id == "" {
		return result, errors.New("Message ID must not be blank")
	}
	
	r, err :=  NewAPIWebRequest(api.config, "message/" + message_id + "/status", "GET", "")
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
			result.StatusCode = APIResultStatusesError;
			result.StatusDescription = "Unable to unmarshal result from API";
		}
		result.MessageStatus = status
		result.StatusCode = APIResultStatusesOk;
		result.StatusDescription = "Ok";
	}
	return result, nil
}

