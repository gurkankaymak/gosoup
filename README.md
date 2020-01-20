# gosoup

[![Go Report Card](https://goreportcard.com/badge/github.com/gurkankaymak/gosoup)](https://goreportcard.com/report/github.com/gurkankaymak/gosoup)

Lightweight go library for pulling data out of HTML

it provides a convenient API for constructing and manipulating data, inspired by BeautifulSoup and Jsoup

## Installation
```go get -u github.com/gurkankaymak/gosoup```

## API Overview
```go
Find(tagName string, attributes Attributes) *Element // returns the first occurrence of the node with the given tagName and attributes
FindByTag(tagName string) *Element // returns the first occurrence of the node with the given tagName
FindByAttributes(attributes Attributes) *Element // returns the first occurrence of the node with the given attributes

FindAll(tagName string, attributes Attributes) []*Element // returns all nodes with the given tagName and attributes
FindAllByTag(tagName string) []*Element // returns all nodes with the given tagName
FindAllByAttributes(attributes Attributes) []*Element // returns all nodes with the given attributes
```

## Usage
```go
element, err := gosoup.ParseAsHTML("... html as string ...")
if err != nil {
    // log/handle error
    return err
}
loginElement := element.Find("form", gosoup.Attributes{"name": "login", "class": "loginTable"})
if loginElement == nil {
    fmt.Println("could not find login element")
}
fmt.Println("loginElement:", loginElement)
```