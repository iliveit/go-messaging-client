//Package messagingapi implements the iliveit Messaging API
package messagingapi

import (
	"fmt"
	"bytes"
	"time"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// NewAPIWebRequest creates a new APIWebRequest instance for the API to use
func NewAPIWebRequest(config APIConfig, url string, method string, data string) (*APIWebRequest, error) {
	
	if url == "" {
		return nil, errors.New("Url can not be blank")
	}
	
	if method == "" {
		return nil, errors.New("Method can not be blank")
	}
	
	r := APIWebRequest{
		config: config,
		Url: url,
		Method: method,
		Data: data,
	}
	
	return &r, nil
}

// Execute the given request using the parameters from NewAPIWebRequest
func (r *APIWebRequest) Execute() (string, int, error) {
	
	timeout := time.Duration(time.Second * 10)
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	
	var jsonBytes = []byte(r.Data)
	req, err := http.NewRequest(r.Method, r.config.Endpoint + r.Url, bytes.NewBuffer(jsonBytes))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + r.config.AccessToken)
	req.ContentLength = int64(len(jsonBytes))
	req.Close = true
	
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			return "", resp.StatusCode, err	
		} else {
			return "", 0, err
		}
	}
	defer resp.Body.Close()

	// Get response body to log error
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	if err != nil {
		return "", resp.StatusCode, err
	}

	if resp.StatusCode == http.StatusOK {
		return bodyString, 200, nil
		
	} else if resp.StatusCode == http.StatusNotFound {
		err := errors.New("Invalid route: 404")
		return "", 404, err
	} else {
		var responseObj WebRequestResponse
		err = json.Unmarshal(body, &responseObj)
		if err != nil {
			err = errors.New(fmt.Sprintf("Unable to submit: %s", bodyString))
			return "", resp.StatusCode, err	
		}
		err = errors.New(responseObj.Error)
		return "", resp.StatusCode, err
	}
	
	return bodyString, resp.StatusCode, nil
}