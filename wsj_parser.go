package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const htm = `<!DOCTYPE html>
<html>
<head>
    <title></title>
</head>
<body>
<div class="article-content">
	<p>more
		<a href="">content</a>
	</p>
    <p>This is here</p>
    <p>Call at </p>
    <p>Hello and </p>
</div>
<div>
<p>Other Paragraph</p>
</body>
</html>`

// simple error handler
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func isNodeArticleDiv(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "class" && attr.Val == "article-content" {
			return true
		}
	}
	return false
}

// TraverseNode collecting the nodes that match the given function
func TraverseNode(doc *html.Node, filter func(node *html.Node) bool) (nodes []*html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if filter(n) {
			fmt.Println(RenderNode(n))
			nodes = append(nodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nodes
}

// find the article content <div> node
func PickArticleNode(doc *html.Node) (article *html.Node) {
	var articleNode *html.Node
	var f func(*html.Node) (article *html.Node)
	f = func(n *html.Node) (article *html.Node) {
		if isNodeArticleDiv(n) {
			articleNode = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if articleNode != nil {
				return articleNode
			}
			f(c)
		}
		if articleNode != nil {
			return articleNode
		}
		return n
	}
	return f(doc)
}

// Works better than: https://github.com/yhat/scrape/blob/master/scrape.go#L129
// because you can cut the search short from the matcher function
func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func test() {

	// data := strings.NewReader(htm)
	data, err := os.Open("sample-site/wsj-article1.html")
	Check(err)
	doc, err := html.Parse(data)
	Check(err)

	filter := func(node *html.Node) (pass bool) {
		if node.DataAtom == atom.P && len(node.Attr) == 0 {
			return true
		}
		return false
	}

	article := PickArticleNode(doc)
	TraverseNode(article, filter)

}
