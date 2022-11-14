package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

// var (
// 	// DefaultHTTPGetAddress Default Address
// 	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

// 	// ErrNoIP No IP found in response
// 	ErrNoIP = errors.New("No IP in HTTP response")

// 	// ErrNon200Response non 200 status code in response
// 	ErrNon200Response = errors.New("Non 200 Response found")
// )

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	path := request.Body
	body := map[string]interface{}{}

	if err := json.Unmarshal([]byte(path), &body); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("invalid body %v", err),
			StatusCode: 400,
		}, nil
	}

	name := fmt.Sprintf("%v", body["name"])
	lastName := fmt.Sprintf("%v", body["lastName"])
	sender := fmt.Sprintf("%v", body["sender"])
	receiver := fmt.Sprintf("%v", body["email"])

	if name == "" || lastName == "" || sender == "" || receiver == "" {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintln("invalid body"),
			StatusCode: 400,
		}, nil
	}

	mySession := session.Must(session.NewSession())
	svc := ses.New(mySession)

	// The subject line for the email.
	Subject := "Hi " + name + lastName + " here is your account summary"

	// The HTML body for the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(receiver),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(EmailBody()),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err email %v", err),
			StatusCode: 200,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Email id %v", *result.MessageId),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)

}
