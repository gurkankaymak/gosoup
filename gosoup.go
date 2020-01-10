package gosoup

import (
	"bytes"
	"encoding/xml"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Element struct {
	*html.Node
}

type Attributes map[string]string

func Html(input string) (*Element, error) {
	err := validateHtml(input)
	if err != nil {
		return nil, errors.New("invalid html")
	}
	rootNode, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return nil, err
	}
	return &Element{Node: rootNode}, nil
}

// html.Parse() does not return error on invalid htmls, so input is validated with encoding/xml
func validateHtml(input string) error {
	decoder := xml.NewDecoder(strings.NewReader(input))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity
	for {
		_, err := decoder.Token()
		switch err {
		case io.EOF:
			return nil
		case nil:
		default:
			return err
		}
	}
}

// Returns the string representation of the parse tree, it panics if there is an error in html.Render()
func (element *Element) String() string {
	var buffer bytes.Buffer
	err := html.Render(&buffer, element.Node)
	panic(err)
	return buffer.String()
}

func (element *Element) Find(tagName string, attributes Attributes) *Element {
	foundElement, _ := find(element.Node, tagName, attributes, true, true)
	return foundElement
}

func (element *Element) FindByTag(tagName string) *Element {
	foundElement, _ := find(element.Node, tagName, nil, true, false)
	return foundElement
}

func (element *Element) FindByAttributes(attributes Attributes) *Element {
	foundElement, _ := find(element.Node, "", attributes, false, true)
	return foundElement
}

func (element *Element) FindAll(tagName string, attributes Attributes) []*Element {
	return findAll(element.Node, tagName, attributes, true, true)
}

func (element *Element) FindAllByTag(tagName string) []*Element {
	return findAll(element.Node, tagName, nil, true, false)
}

func (element *Element) FindAllByAttributes(attributes Attributes) []*Element {
	return findAll(element.Node, "", attributes, false, true)
}

func (element *Element) getAttribute(attributeName string) (value string, ok bool) {
	for _, attribute := range element.Attr {
		if attribute.Key == attributeName {
			return attribute.Val, true
		}
	}
	return "", false
}

func find(node *html.Node, tagName string, attributes Attributes, includeTagName, includeAttributes bool) (*Element, bool) {
	if node.Type == html.ElementNode {
		found := checkNode(node, tagName, attributes, includeTagName, includeAttributes)
		if found {
			return &Element{node}, found
		}
	}
	// check child nodes
	for nextNode := node.FirstChild; nextNode != nil; nextNode = nextNode.NextSibling {
		element, found := find(nextNode, tagName, attributes, includeTagName, includeAttributes)
		if found {
			return element, found
		}
	}
	return nil, false
}

func findAll(node *html.Node, tagName string, attributes Attributes, includeTagName, includeAttributes bool) []*Element {
	var foundElements []*Element
	if node.Type == html.ElementNode {
		found := checkNode(node, tagName, attributes, includeTagName, includeAttributes)
		if found {
			foundElements = append(foundElements, &Element{node})
		}
	}
	// check child nodes
	for nextNode := node.FirstChild; nextNode != nil; nextNode = nextNode.NextSibling {
		elements := findAll(nextNode, tagName, attributes, includeTagName, includeAttributes)
		foundElements = append(foundElements, elements...)
	}
	return foundElements
}

func checkNode(node *html.Node, tagName string, attributes Attributes, includeTagName, includeAttributes bool) bool {
	foundByTag := includeTagName && node.Data == tagName
	foundByAttributes := includeAttributes && checkAttributes(node, attributes)
	if includeTagName && includeAttributes {
		return foundByTag && foundByAttributes
	} else {
		return foundByTag || foundByAttributes
	}
}

func checkAttributes(node *html.Node, attributes Attributes) bool {
	var found bool
	for name, value := range attributes {
		found = false
		for _, attribute := range node.Attr {
			if attribute.Key == name && attribute.Val == value {
				found = true
				break
			}
		}
		if !found { // if one of the given attributes is not found, no need to look for others
			return false
		}
	}
	return true
}
