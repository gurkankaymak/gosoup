package gosoup

import (
	"testing"
)

const testHtml = `
	<div>
		<div id="divId1" class="class1">
		
		</div>
		<div id="divId2" class="class1">
		
		</div>
		<div id="divId3" class="class2">
		
		</div>
	</div>
	`
var rootElement, _ = Html(testHtml)

func TestHtml(t *testing.T) {
	t.Run("should return error if the input is not a valid html", func(t *testing.T) {
		testHtml := "<div></"
		_, err := Html(testHtml)
		if err == nil {
			t.Errorf("did not get the expected invalid html error")
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("should return nil if the element with the given tagName and attributes does not exist", func(t *testing.T) {
		element := rootElement.Find("div", Attributes{"id": "testDiv5"})
		if element != nil {
			t.Errorf("returned element: %q, expected: 'nil'", element)
		}
	})

	t.Run("should find the first element with the given tagName and attributes", func(t *testing.T) {
		expectedTag := "div"
		expectedAttrKey := "id"
		expectedAttrVal := "divId1"
		element := rootElement.Find(expectedTag, Attributes{expectedAttrKey: expectedAttrVal})
		if element == nil {
			t.Fatal("could not find expected element")
		}
		if element.Data != expectedTag {
			t.Errorf("wrong element tag, expected: %q, actual: %q", expectedTag, element.Data)
		}
		attributeValue, ok := element.getAttribute(expectedAttrKey)
		if !ok || attributeValue != expectedAttrVal {
			t.Errorf("expected attribute: %q: %q does not exist", expectedAttrKey, expectedAttrVal)
		}
	})
}

func TestFindAll(t *testing.T) {
	t.Run("should find all elements with the given tagName and attributes", func(t *testing.T) {
		expectedTag := "div"
		expectedAttrKey := "class"
		expectedAttrVal := "class1"
		expectedElementsSize := 2
		elements := rootElement.FindAll(expectedTag, Attributes{expectedAttrKey: expectedAttrVal})
		if len(elements) != expectedElementsSize {
			t.Errorf("wrong number of elements found: %q, expected number: %q", len(elements), expectedElementsSize)
		}
		for _, element := range elements {
			if element.Data != expectedTag {
				t.Errorf("wrong element tag, expected: %q, actual: %q", expectedTag, element.Data)
			}
			attributeValue, ok := element.getAttribute(expectedAttrKey)
			if !ok || attributeValue != expectedAttrVal {
				t.Errorf("expected attribute: %q: %q does not exist", expectedAttrKey, expectedAttrVal)
			}
		}
	})
}


