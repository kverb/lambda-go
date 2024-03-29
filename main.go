package main

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	data := strings.NewReader(request.Body)
	doc, err := html.Parse(data)
	Check(err)
	filter := func(node *html.Node) (pass bool) {
		if node.DataAtom == atom.P && len(node.Attr) == 0 {
			return true
		}
		return false
	}

	article := PickArticleNode(doc)
	nodes := TraverseNode(article, filter)

	bodyParts := make([]string, 0)
	for _, node := range nodes {
		bodyParts = append(bodyParts, RenderNode(node))
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       strings.Join(bodyParts[:], ""),
		Headers:    headers,
	}, nil
}

func main() {
	// test()
	// TestHandler()
	lambda.Start(handler)
}
