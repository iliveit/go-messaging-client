// Package main is a sample application for the Go Messaging API client library
package main

import(
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/iliveit/go-messaging-client"
)

var api *messagingapi.MessagingAPI

func main() {
	fmt.Println("Sample Go App for Messaging API")
	
	apiConfig := messagingapi.APIConfig{
		Endpoint: "http://127.0.0.1:9000",
		AccessToken: "testtoken",
	}
	
	var err error
	api, err = messagingapi.New(apiConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} 
	
	//SamplePing()
	//SampleSubmitSMS()
	//SampleSubmitMMS()
	//SampleSubmitEmail()
	
	// Sample server to handle incoming SMS
	//http.HandleFunc("/", HandleIncomingSMS)
	//http.ListenAndServe(":9001", nil)
	
}

// SamplePing shows how to use the client to Ping the API. The Ping
// route checks the access token, so this is regarded as the
// easiest way to verify your application correctly authenticates
func SamplePing() {
	result, err := api.Ping()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	if result.StatusCode != messagingapi.APIResultStatusesOk {
		fmt.Println(result.StatusDescription)
	} else {
		fmt.Println("Success")
	}
}

// SampleSubmitSMS shows how to submit an SMS message using the client library
func SampleSubmitSMS() {
	// Create the wrapper object
	msg := messagingapi.NewMessage{
		Action: messagingapi.APIActionTypesSubmitSMS,
		MVNOID: 2,
		Campaign: "GoClientTest",
		PostbackReplyUrl: "http://127.0.0.1:9001",
	}
	// Create the message data
	msgData := messagingapi.SubmitSMSMessageData{
		Network: "local_smpp",
		MSISDN: []string{"27760913077"},
		Message: "This is my SMS text",
		ExtraDigits: "00015",
	}
	msg.Data = msgData
	// Send the create request
	result, err := api.Create(msg)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	// Handle the result
	if result.StatusCode != messagingapi.APIResultStatusesOk {
		fmt.Println(result.StatusDescription)
	} else {
		fmt.Println("Success")
		fmt.Println(result.MessageResult.MessageID)
	}
}

// SampleSubmitMMS shows how to submit an MMS message using the client library
func SampleSubmitMMS() {
	// Create the wrapper object
	msg := messagingapi.NewMessage{
		Action: messagingapi.APIActionTypesSubmitMMS,
		MVNOID: 2,
		Campaign: "GoClientTest",
	}
	// Create the message data
	msgData := messagingapi.SubmitMMSMessageData{
		// Network '*' uses the portability list to determine
		// the destination network
		Network: "*",
		MSISDN: []string{"27760913077"},
		Subject: "MMS Subject",
	}
	
	slideContent := messagingapi.MMSSlideContent{
		Type: messagingapi.MMSContentTypeText,
		Mime: "text/plain",
        Name: "TextDocument1.txt",
        Data: "My Plain Text MMS",
	}
	
	slide := messagingapi.MMSSlide{
		Duration: "10",
	}
	slide.Content = append(slide.Content, slideContent)
	msgData.Slides = append(msgData.Slides, slide)
	
	
	msg.Data = msgData
	// Send the create request
	result, err := api.Create(msg)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	// Handle the result
	if result.StatusCode != messagingapi.APIResultStatusesOk {
		fmt.Println(result.StatusDescription)
	} else {
		fmt.Println("Success")
		fmt.Println(result.MessageResult.MessageID)
	}
}

// SampleSubmitEmail shows how to submit an Email message using the client library
func SampleSubmitEmail() {
	// Create the wrapper object
	msg := messagingapi.NewMessage{
		Action: messagingapi.APIActionTypesSubmitEmail,
		MVNOID: 2,
		Campaign: "GoClientTest",
	}
	// Create the message data
	msgData := messagingapi.SubmitEmailMessageData{
		Network: "local_email",
		Address: []string{"donovan.solms@gmail.com"},
		Subject: "Email Subject",
		HTML: "<h1>This is my email</h1>",
		Text: "This is plain text part",
	}
	
	msg.Data = msgData
	// Send the create request
	result, err := api.Create(msg)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	// Handle the result
	if result.StatusCode != messagingapi.APIResultStatusesOk {
		fmt.Println(result.StatusDescription)
	} else {
		fmt.Println("Success")
		fmt.Println(result.MessageResult.MessageID)
	}
}

// HandleIncomingSMS receives POSTs from the API for incoming SMS
func HandleIncomingSMS(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	
    var incoming messagingapi.IncomingSMS
    err := decoder.Decode(&incoming)
    if err != nil {
        fmt.Println("Unable to get POST body: " + err.Error())
        return
    }
    
	fmt.Println("MessageID: " + incoming.MessageId)
	fmt.Println("SourceMSISDN: " + incoming.SourceMSISDN)
	fmt.Println("DestinationMSISDN: " + incoming.DestinationMSISDN)
	fmt.Println("Message: " + incoming.Message)
	fmt.Println("ExtraDigits: " + incoming.ExtraDigits)
	
	// You must respond with a status code of 200 otherwise the API will keep retrying
	// the message to this endpoint
	w.Write([]byte("Ok"))
}