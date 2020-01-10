package gosoup

import (
	"bytes"
	"encoding/xml"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Element is the root node returned from the api methods, it contains a pointer to the html.Node
// the underlying *html.Node could be accessed as element.Node
type Element struct {
	*html.Node
}

// Attributes is a map[string]string that represents the attributes of an element
// it is used to call the Find and derivative methods conveniently as in the BeautifulSoup
// in example: Find("div", Attributes{"class":"exampleClass", "name":"exampleName"})
type Attributes map[string]string

// ParseAsHTML function Validates and parses the given string as html
// Returns an Element pointer to the root node and error if any error occurs
func ParseAsHTML(input string) (*Element, error) {
	err := validateHTML(input)
	if err != nil {
		return nil, errors.New("invalid html: " + input)
	}
	rootNode, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return nil, err
	}
	return &Element{Node: rootNode}, nil
}

// html.Parse() does not return error on invalid htmls, so input is validated with encoding/xml
func validateHTML(input string) error {
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

// Find method returns the first occurrence of the node with the given tagName and attributes
// returns nil if not found any node with the given parameters
func (element *Element) Find(tagName string, attributes Attributes) *Element {
	foundElement, _ := find(element.Node, tagName, attributes, true, true)
	return foundElement
}

// FindByTag method returns the first occurrence of the node with the given tagName, returns nil if not found
func (element *Element) FindByTag(tagName string) *Element {
	foundElement, _ := find(element.Node, tagName, nil, true, false)
	return foundElement
}

// FindByAttributes methods returns the first occurrence of the node with the given attributes
func (element *Element) FindByAttributes(attributes Attributes) *Element {
	foundElement, _ := find(element.Node, "", attributes, false, true)
	return foundElement
}

// FindAll method returns all occurrences of the nodes with the given tagName and attributes
func (element *Element) FindAll(tagName string, attributes Attributes) []*Element {
	return findAll(element.Node, tagName, attributes, true, true)
}

// FindAllByTag method returns all occurrences of the nodes with the given tagName
func (element *Element) FindAllByTag(tagName string) []*Element {
	return findAll(element.Node, tagName, nil, true, false)
}

// FindAllByAttributes method returns all occurrences of the nodes with the given attributes
func (element *Element) FindAllByAttributes(attributes Attributes) []*Element {
	return findAll(element.Node, "", attributes, false, true)
}

// GetAttribute method behaves like a map lookup, returns the value of the attribute with the given key and true if found,
// returns "" and false if not found
func (element *Element) GetAttribute(attributeName string) (string, bool) {
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
	}
	return foundByTag || foundByAttributes
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
