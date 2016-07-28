// Package main is a sample application for the Go Messaging API client library
package main

import (
	"fmt"
	"os"
	//"time"
	"encoding/json"
	"github.com/iliveit/go-messaging-client"
	"net/http"
)

var api *messagingapi.MessagingAPI

func main() {
	fmt.Println("Sample Go App for Messaging API")

	apiConfig := messagingapi.APIConfig{
		Endpoint:    "api endpoint",
		AccessToken: "your access token",
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
	//SampleGetMessageStatus()
	//SampleSubmitPaymentReminder()

	// You need to specify the type of message for the approval batch
	batchId := SampleCreateApprovalBatch(messagingapi.APIActionTypesSubmitEmail)
	// Send some messages
	SampleSubmitEmailWithApproval(batchId, "12345")
	SampleSubmitEmailWithApproval(batchId, "abcde")
	SampleSubmitEmailWithApproval(batchId, "1a2b3")
	//SampleSubmitSMSWithApproval(batchId)
	//SampleSubmitMMSWithApproval(batchId)
	//SampleSubmitStatementWithApproval(batchId)
	SampleApprovalUpdate(batchId)

	// Sample server to handle incoming SMS
	http.HandleFunc("/", HandleIncomingSMS)
	http.HandleFunc("/status", HandleStatusUpdates)
	http.ListenAndServe(":9001", nil)

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
		Action:              messagingapi.APIActionTypesSubmitSMS,
		MVNOID:              4,
		Campaign:            "GoClientTest",
		PostbackReplyUrl:    "http://127.0.0.1:9001",
		PostbackStatusUrl:   "http://127.0.0.1:9001/status",
		PostbackStatusTypes: "submit,archive,sent,delivery",
	}
	// Create the message data
	msgData := messagingapi.SubmitSMSMessageData{
		Network:     "local_smpp",
		MSISDN:      []string{"277777"},
		Message:     "This is my SMS text",
		ExtraDigits: "00123",
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
		fmt.Println("Sent with ExtraDigits: %s", msgData.ExtraDigits)
	}
}

// SampleSubmitMMS shows how to submit an MMS message using the client library
func SampleSubmitMMS() {
	// Create the wrapper object
	msg := messagingapi.NewMessage{
		Action:              messagingapi.APIActionTypesSubmitMMS,
		MVNOID:              2,
		Campaign:            "GoClientTest",
		PostbackStatusUrl:   "http://127.0.0.1:9001/status",
		PostbackStatusTypes: "submit,sent,delivery,archive",
	}
	// Create the message data
	msgData := messagingapi.SubmitMMSMessageData{
		// Network '*' uses the portability list to determine
		// the destination network
		Network: "*",
		MSISDN:  []string{"27777777"},
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
		Action:              messagingapi.APIActionTypesSubmitEmail,
		MVNOID:              2,
		Campaign:            "GoClientTest",
		PostbackStatusUrl:   "http://127.0.0.1:9001/status",
		PostbackStatusTypes: "build,submit,sent,delivery,archive",
	}
	// Create the message data
	msgData := messagingapi.SubmitEmailMessageData{
		Network: "local_email",
		Address: []string{"none@example.com"},
		Subject: "Email Subject",
		HTML:    "<h1>This is my email</h1>",
		Text:    "This is plain text part",
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

// SampleGetMessageStatus shows how to retrieve the status of a message
func SampleGetMessageStatus() {

	result, err := api.GetMessageStatus("1-90aaf0b4-65d2-4556-7ec1-a134afdf6e76")
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	// Handle the result
	if result.StatusCode != messagingapi.APIResultStatusesOk {
		fmt.Println(result.StatusDescription)
	} else {
		fmt.Println("Success")
		fmt.Println(result.MessageStatus.Campaign)
		fmt.Println(result.MessageStatus.BuildStatus)
		fmt.Println(result.MessageStatus)
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

	byteString, _ := json.Marshal(incoming)
	fmt.Println(string(byteString))

	fmt.Println("MessageID: " + incoming.MessageId)
	fmt.Println("SourceMSISDN: " + incoming.SourceMSISDN)
	fmt.Println("DestinationMSISDN: " + incoming.DestinationMSISDN)
	fmt.Println("Message: " + incoming.Message)
	fmt.Println("ExtraDigits: " + incoming.ExtraDigits)

	// You must respond with a status code of 200 otherwise the API will keep retrying
	// the message to this endpoint
	w.Write([]byte("Ok"))
}

// SampleCreateApprovalBatch creates a simple
// approval batch before loading data
func SampleCreateApprovalBatch(messageActionType uint32) uint32 {
	fmt.Println("Create an approval batch")

	var internalApprovalPeople []messagingapi.ApprovalPerson
	internalApprovalPeople = append(internalApprovalPeople, messagingapi.ApprovalPerson{
		Name:  "John Internal",
		Email: "none@example.com",
	})
	var externalApprovalPeople []messagingapi.ApprovalPerson
	externalApprovalPeople = append(externalApprovalPeople, messagingapi.ApprovalPerson{
		Name:  "John External",
		Email: "none@example.com",
	})

	approvalRequest := messagingapi.ApprovalRequest{
		ActionType:     messageActionType,
		Name:           "Demo Approval Batch",
		MVNOID:         1,
		MaxApprovals:   10,
		InternalPeople: internalApprovalPeople,
		ExternalPeople: externalApprovalPeople,
	}

	result, err := api.CreateApproval(approvalRequest)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	// Handle the result
	if result.StatusCode != messagingapi.APIResultStatusesOk {
		fmt.Println(result.StatusDescription)
	} else {
		fmt.Println("Success")
		fmt.Println(result.RequestResult.BatchID)
		return result.RequestResult.BatchID
	}
	return 0
}

// SampleCreateApprovalBatch creates a simple
// approval batch before loading data
func SampleApprovalUpdate(batchId uint32) {
	fmt.Println("Update an approval batch with report")

	// Build the report as CSV line
	var reports []messagingapi.CsvReport

	var reportDataA []string
	reportDataA = append(reportDataA, "this,is,a,line")
	reportDataA = append(reportDataA, "what,is,my,time")
	reports = append(reports, messagingapi.CsvReport{
		Filename: "Report A.csv",
		Lines:    reportDataA,
	})

	var reportDataB []string
	reportDataB = append(reportDataB, "Total,Col A,Col B")
	reportDataB = append(reportDataB, "100 000,95 000,5 000")
	reports = append(reports, messagingapi.CsvReport{
		Filename: "Report B.csv",
		Lines:    reportDataB,
	})

	updateRequest := messagingapi.ApprovalUpdateRequest{
		BatchID: batchId,
		State:   messagingapi.ApprovalBatchStateDataReceived,
		Reports:  reports,
	}

	result, err := api.UpdateApproval(updateRequest)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	// Handle the result
	if result.StatusCode != messagingapi.APIResultStatusesOk {
		fmt.Println(result.StatusDescription)
	} else {
		fmt.Println("Success")
		fmt.Println(result.RequestResult.BatchID)
	}
}

// SampleSubmitPaymentReminder sends a payment reminder
func SampleSubmitPaymentReminder() {

	fmt.Println("Sending Payment Reminder")

	buildRequest := messagingapi.BuildRequest{
		MVNOID:           2,
		Campaign:         "GoTest",
		BuildTemplate:    15,
		AfterBuildAction: messagingapi.APIActionTypesSubmitMMS,
	}
	buildRequest.Data = "{\"CustomerName\":\"John Doe\",\"AccountNumber\":\"AC0001\",\"AmountDue\":100.00}"
	msgData := messagingapi.SubmitMMSMessageData{
		// Network '*' uses the portability list to determine
		// the destination network
		Network: "*",
		MSISDN:  []string{"270000000"},
	}
	buildRequest.AfterBuildData = msgData

	result, err := api.Generate(buildRequest)
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
func SampleSubmitEmailWithApproval(batchId uint32, unique string) {
	// Create the wrapper object
	msg := messagingapi.NewMessage{
		Action:        messagingapi.APIActionTypesSubmitEmail,
		MVNOID:        2,
		ApprovalBatch: batchId,
		Campaign:      "ApprovalTest",
		//PostbackStatusUrl:   "http://127.0.0.1:9001/status",
		//PostbackStatusTypes: "build,submit,sent,delivery,archive",
	}

	var attachments []messagingapi.EmailAttachment
	attachment := messagingapi.EmailAttachment{
		Filename: "TestDocument.pdf",
		Data:     "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MCBSCi9Qcm9jU2V0Wy9QREYvVGV4dF0KPj4KZW5kb2JqCgoxIDAgb2JqCjw8L1R5cGUvUGFnZS9QYXJlbnQgNCAwIFIvUmVzb3VyY2VzIDExIDAgUi9NZWRpYUJveFswIDAgNTk1IDg0Ml0vR3JvdXA8PC9TL1RyYW5zcGFyZW5jeS9DUy9EZXZpY2VSR0IvSSB0cnVlPj4vQ29udGVudHMgMiAwIFI+PgplbmRvYmoKCjQgMCBvYmoKPDwvVHlwZS9QYWdlcwovUmVzb3VyY2VzIDExIDAgUgovTWVkaWFCb3hbIDAgMCA1OTUgODQyIF0KL0tpZHNbIDEgMCBSIF0KL0NvdW50IDE+PgplbmRvYmoKCjEyIDAgb2JqCjw8L1R5cGUvQ2F0YWxvZy9QYWdlcyA0IDAgUgovT3BlbkFjdGlvblsxIDAgUiAvWFlaIG51bGwgbnVsbCAwXQovTGFuZyhlbi1aQSkKPj4KZW5kb2JqCgoxMyAwIG9iago8PC9DcmVhdG9yPEZFRkYwMDU3MDA3MjAwNjkwMDc0MDA2NTAwNzI+Ci9Qcm9kdWNlcjxGRUZGMDA0QzAwNjkwMDYyMDA3MjAwNjUwMDRGMDA2NjAwNjYwMDY5MDA2MzAwNjUwMDIwMDAzNTAwMkUwMDMxPgovQ3JlYXRpb25EYXRlKEQ6MjAxNjA1MjQxMDA2NTErMDInMDAnKT4+CmVuZG9iagoKeHJlZgowIDE0CjAwMDAwMDAwMDAgNjU1MzUgZiAKMDAwMDAwNzIzNSAwMDAwMCBuIAowMDAwMDAwMDE5IDAwMDAwIG4gCjAwMDAwMDAyMjYgMDAwMDAgbiAKMDAwMDAwNzM3OCAwMDAwMCBuIAowMDAwMDAwMjQ2IDAwMDAwIG4gCjAwMDAwMDYzNjkgMDAwMDAgbiAKMDAwMDAwNjM5MCAwMDAwMCBuIAowMDAwMDA2NTg1IDAwMDAwIG4gCjAwMDAwMDY5MzYgMDAwMDAgbiAKMDAwMDAwNzE0OCAwMDAwMCBuIAowMDAwMDA3MTgwIDAwMDAwIG4gCjAwMDAwMDc0NzcgMDAwMDAgbiAKMDAwMDAwNzU3NCAwMDAwMCBuIAp0cmFpbGVyCjw8L1NpemUgMTQvUm9vdCAxMiAwIFIKL0luZm8gMTMgMCBSCi9JRCBbIDwxNkQ0OEQwMzcyN0UwNkIwMzNCQ0M4MjRFMjNFQTgwOD4KPDE2RDQ4RDAzNzI3RTA2QjAzM0JDQzgyNEUyM0VBODA4PiBdCi9Eb2NDaGVja3N1bSAvRTZBOTE0QjJENDUzQTQwNjQ0OTlFQjJDRkQ5MjAxNjUKPj4Kc3RhcnR4cmVmCjc3NDkKJSVFT0YK",
	}
	attachments = append(attachments, attachment)
	// Create the message data

	msgData := messagingapi.SubmitEmailMessageData{
		Network:     "local_email",
		Address:     []string{"none@example.com"},
		Subject:     "Email Subject " + unique,
		HTML:        "<h1>This is my email " + unique + "</h1>",
		Text:        "This is plain text part " + unique,
		Attachments: attachments,
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

// SampleSubmitMTNWithApproval submits a build request with approval request
func SampleSubmitStatementWithApproval(batchId uint32) {

	fmt.Println("Sending build request")

	buildRequest := messagingapi.BuildRequest{
		MVNOID:           1,
		Campaign:         "Statement Sample",
		BuildTemplate:    2,
		AfterBuildAction: messagingapi.APIActionTypesSubmitMMS,
		ApprovalBatch:    batchId,
	}
	buildRequest.Data = "{\"AccountNumber\":\"123432\",\"Address\":[\"123 My Street\",\"MyCity\"],\"AmountDue\":10.0,\"BankAccountNumber\":\"654323456\",\"BankName\":\"Absa Bank\",\"BranchCode\":\"123654\",\"ClosingBalance\":9.0,\"CurrentBalance\":4.0,\"MSISDN\":\"27000\",\"Name\":\"John Doe\",\"OpeningBalance\":14.0,\"PaymentDue\":\"2015-05-29\",\"PaymentType\":\"Cash\",\"PostalCode\":\"0181\",\"TotalOutstandingBalance\":325.0,\"Transactions\":[{\"Amount\":123.4,\"Description\":\"Only invoice line\"}],\"VatNo\":\"111122233445\",\"ThirtyDaysOverdue\":10.0,\"ThirtyDaysOverdueText\":\"30 Days Overdue\",\"SixtyDaysOverdue\":28.0,\"SixtyDaysOverdueText\":\"60 Days Overdue\",\"NinetyDaysOverdue\":0.0,\"NinetyDaysOverdueText\":\"90 Days Overdue\",\"Distribution\":null}"
	msgData := messagingapi.SubmitMMSMessageData{
		// Network '*' uses the portability list to determine
		// the destination network
		Network: "mtn",
		MSISDN:  []string{"27000"},
	}
	buildRequest.AfterBuildData = msgData

	result, err := api.Generate(buildRequest)
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

// HandleIncomingSMS receives POSTs from the API for status updates
func HandleStatusUpdates(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var status messagingapi.StatusResult
	err := decoder.Decode(&status)
	if err != nil {
		fmt.Println("Unable to get POST body: " + err.Error())
		return
	}

	byteString, _ := json.MarshalIndent(status, "", " ")
	fmt.Println(string(byteString))
	fmt.Println("\n\n\n")

	/*
		fmt.Println("MessageID: " + incoming.MessageId)
		fmt.Println("SourceMSISDN: " + incoming.SourceMSISDN)
		fmt.Println("DestinationMSISDN: " + incoming.DestinationMSISDN)
		fmt.Println("Message: " + incoming.Message)
		fmt.Println("ExtraDigits: " + incoming.ExtraDigits)
	*/

	// You must respond with a status code of 200 otherwise the API will keep retrying
	// the message to this endpoint
	w.Write([]byte("Ok"))
}
