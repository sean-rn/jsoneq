package jsoneq

import (
	"encoding/json"
	"fmt"
)

// TestingT is an interface wrapper around the methods of *testing.T that we require
type TestingT interface {
	Errorf(format string, args ...any)
}

// JSONEq asserts that two JSON strings are equivalent, ignoring the order of array elements.
//
//	jsoneq.JSONEq(t, `["a", "b", 1, 2]`, `[2, "b", 1, "a"]`)
func JSONEq(t TestingT, expected string, actual string, msgAndArgs ...any) bool {
	if h, ok := t.(interface{ Helper() }); ok {
		h.Helper()
	}
	var expectedJSONAsInterface, actualJSONAsInterface any

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		t.Errorf("Expected ('%s') is not valid json.\n%s\n%s", expected, err.Error(), messageFromMsgAndArgs(msgAndArgs...))
		return false
	}
	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		t.Errorf("Actual ('%s') is not valid json.\n%s\n%s", actual, err.Error(), messageFromMsgAndArgs(msgAndArgs...))
	}

	if compareJSONThings(expectedJSONAsInterface, actualJSONAsInterface) {
		return true // They're equal, great!
	}

	// Oh, they aren't so equal.  Ah well
	t.Errorf("Not equivalent JSON:\nexpected: %s\nactual:   %s\n%s", expected, actual, messageFromMsgAndArgs(msgAndArgs...))
	return false
}

// We assume that a or b can ONLY be the following things, since that is what JSON marshal does:
//   - bool, for JSON booleans
//   - float64, for JSON numbers
//   - string, for JSON strings
//   - []any, for JSON arrays
//   - map[string]any, for JSON objects
//   - nil for JSON null
func compareJSONThings(a, b any) bool {
	switch aval := a.(type) {
	case nil:
		return a == b
	case bool:
		if bval, ok := b.(bool); ok {
			return aval == bval
		}
	case float64:
		if bval, ok := b.(float64); ok {
			return aval == bval
		}
	case string:
		if bval, ok := b.(string); ok {
			return aval == bval
		}
	case []any:
		if bval, ok := b.([]any); ok {
			return compareJSONArrays(aval, bval)
		}
	case map[string]any:
		if bval, ok := b.(map[string]any); ok {
			return compareJSONMaps(aval, bval)
		}
	default:
		panic(fmt.Errorf("invalid type %T of %#v", a, a))
	}
	return false
}

func compareJSONMaps(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false // Shortcut
	}

	for key, aval := range a {
		bval, exists := b[key]
		if !exists {
			return false // Key present in A but not B
		} else if !compareJSONThings(aval, bval) {
			return false // Values differ
		}
	}

	for key := range b {
		if _, exists := a[key]; !exists {
			return false // Key present in B but not A
		}
	}

	return true
}

func compareJSONArrays(aList, bList []any) bool {
	var extraA, extraB []any
	aLen := len(aList)
	bLen := len(bList)

	// Mark indexes in bValue that we already used
	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := aList[i]
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if compareJSONThings(bList[j], element) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, bList[j])
	}

	if len(extraA) == 0 && len(extraB) == 0 {
		return true
	}
	return false
}

// Format messageFromMsgAndArgs like testify does
func messageFromMsgAndArgs(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}
