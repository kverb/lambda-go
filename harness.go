package main

import (
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler() {
	var request events.APIGatewayProxyRequest
	data, err := ioutil.ReadFile("sample-site/wsj-article1.html")
	Check(err)
	request.Body = string(data)
	resp, err := handler(request)
	Check(err)
	log.Print(resp.Body)
}
