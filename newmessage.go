package messagingapi

import (
	"errors"
)

// Validate checks that all required fields are set before submitting
func (message *NewMessage) Validate() error {
	
	var err error
	if message.Action == 0 {
		err = errors.New("Action must be set")
	}
	if message.MVNOID == 0 {
		err = errors.New("MVNOID must be set and not zero")
	}
	
	if message.Action == APIActionTypesSubmitMMS {
		data, ok := message.Data.(SubmitMMSMessageData)
		if ok {
			if data.Network == "" {
				err = errors.New("Network cannot be blank")
			}
			if len(data.Slides) == 0 {
				err = errors.New("MMS messages must have at least one slide")
			}
			if data.Subject == "" {
				err = errors.New("MMS message must have a subject set")
			}
			if len(data.MSISDN) == 0  || len(data.MSISDN) > 1 {
				err = errors.New("A message must have one recipient set in MSISDN")
			}
		}
		
	} else if message.Action == APIActionTypesSubmitSMS {
		data, ok := message.Data.(SubmitSMSMessageData)
		if ok {
			if data.Network == "" {
				err = errors.New("Network cannot be blank")
			}
			if len(data.MSISDN) == 0  || len(data.MSISDN) > 1 {
				err = errors.New("A message must have one recipient set in MSISDN")
			}
			if data.Message == "" {
				err = errors.New("Message text cannot be blank")
			}
		}
		
	} else if message.Action == APIActionTypesSubmitEmail {
		data, ok := message.Data.(SubmitEmailMessageData)
		if ok {
			if data.Network == "" {
				err = errors.New("Network cannot be blank")
			}
			if len(data.Address) == 0  || len(data.Address) > 1 {
				err = errors.New("Email messages must have at least one recipient listed in Address")
			}
			if data.HTML == "" && data.Text == "" {
				err = errors.New("Email messages must have either HTML or Text set, or both")
			}
			if data.Subject == "" {
				err = errors.New("Email messages must have a subject")
			}
		}
		
	} 
	
	return err
}