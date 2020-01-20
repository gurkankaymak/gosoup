package gosoup

import (
	"testing"
)

const testHTML =
	`<div name="outerDiv">
		<div id="divId1" class="class1">
		
		</div>
		<div id="divId2" class="class1">
		
		</div>
		<div id="divId3" class="class2">
		
		</div>
	</div>`

var rootElement, _ = ParseAsHTML(testHTML)

func TestParseAsHTML(t *testing.T) {
	t.Run("should not return an error or nil element when the given html string is valid", func(t *testing.T) {
		element, err := ParseAsHTML(testHTML)
		if err != nil {
			t.Fatalf("Unexpected error returned, err: %q", err)
		}
		if element == nil {
			t.Errorf("element should not be 'nil' for a valid HTML")
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
		attributeValue, ok := element.GetAttribute(expectedAttrKey)
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
			attributeValue, ok := element.GetAttribute(expectedAttrKey)
			if !ok || attributeValue != expectedAttrVal {
				t.Errorf("expected attribute: %q: %q does not exist", expectedAttrKey, expectedAttrVal)
			}
		}
	})
}

func TestGetAttribute(t *testing.T) {
	t.Run("should return the attribute of element with the given name", func(t *testing.T) {
		firstDivElement := rootElement.FindByTag("div")
		attribute, found := firstDivElement.GetAttribute("name")
		if !found {
			t.Fatal("could not find name attribute")
		}
		if attribute == "" {
			t.Errorf("empty attribute value received, expected: 'outerDiv'")
		}
	})
	t.Run("return empty string and false for the 'found' if the element with the given name does not exist", func(t *testing.T) {
		firstDivElement := rootElement.FindByTag("div")
		attribute, found := firstDivElement.GetAttribute("class")
		if found {
			t.Fatalf("wrong found value: %t, expected: 'false'", found)
		}
		if attribute != "" {
			t.Errorf("wrong attribute value: %q, expected empty string", attribute)
		}
	})
}

func TestString(t *testing.T) {
	firstDivElement := rootElement.FindByTag("div")
	str := firstDivElement.String()
	if str != testHTML {
		t.Errorf("wrong string returned: %q, expected: %q", str, testHTML)
	}
}
