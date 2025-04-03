package jsoneq_test

import (
	"testing"

	. "github.com/sean-rn/jsoneq"
)

func Test_Example(t *testing.T) {
	JSONEq(t, `["a", "b", 1, 2]`, `[2, "b", 1, "a"]`)
}

func TestJSONEq_EqualSONString(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		`{"hello": "world", "foo": "bar"}`,
		`{"hello": "world", "foo": "bar"}`)
	if retVal != true {
		t.Errorf("Expected %v but got %v", true, retVal)
	}
}

func TestJSONEq_EquivalentButNotEqual(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		`{"hello": "world", "foo": "bar"}`,
		`{"foo": "bar", "hello": "world"}`)
	if retVal != true {
		t.Errorf("Expected %v but got %v", true, retVal)
	}
}

func TestJSONEq_HashOfArraysAndHashes(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT, `{
			"numeric": 1.5,
			"array": [{"foo": "bar"}, 1, "string", ["nested", "array", 5.5]],
			"hash": {"nested": "hash", "nested_slice": ["this", "is", "nested"]},
			"string": "foo"
		}`, `{
			"numeric": 1.5,
			"hash": {"nested": "hash", "nested_slice": ["this", "is", "nested"]},
			"string": "foo",
			"array": [{"foo": "bar"}, 1, "string", ["nested", "array", 5.5]]
		}`)
	if retVal != true {
		t.Errorf("Expected %v but got %v", true, retVal)
	}
}

func TestJSONEq_Array(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		`["foo", {"hello": "world", "nested": "hash"}]`,
		`["foo", {"nested": "hash", "hello": "world"}]`)
	if retVal != true {
		t.Errorf("Expected %v but got %v", true, retVal)
	}
}

func TestJSONEq_HashAndArrayNotEquivalent(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		`["foo", {"hello": "world", "nested": "hash"}]`,
		`{"foo": "bar", {"nested": "hash", "hello": "world"}}`)
	if retVal != false {
		t.Errorf("Expected %v but got %v", false, retVal)
	}
}

func TestJSONEq_HashesNotEquivalent(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		`{"foo": "bar"}`,
		`{"foo": "bar", "hello": "world"}`)
	if retVal != false {
		t.Errorf("Expected %v but got %v", false, retVal)
	}
}

func TestJSONEq_ActualIsNotJSON(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		`{"foo": "bar"}`,
		"Not JSON")
	if retVal != false {
		t.Errorf("Expected %v but got %v", false, retVal)
	}
}

func TestJSONEq_ExpectedIsNotJSON(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		"Not JSON",
		`{"foo": "bar", "hello": "world"}`)
	if retVal != false {
		t.Errorf("Expected %v but got %v", false, retVal)
	}
}

func TestJSONEq_ExpectedAndActualNotJSON(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		"Not JSON",
		"Not JSON")
	if retVal != false {
		t.Errorf("Expected %v but got %v", false, retVal)
	}
}

func TestJSONEq_ArraysOfDifferentOrder(t *testing.T) {
	mockT := new(testing.T)
	retVal := JSONEq(mockT,
		`["foo", {"hello": "world", "nested": "hash"}]`,
		`[{ "hello": "world", "nested": "hash"}, "foo"]`)
	if retVal != true {
		t.Errorf("Expected %v but got %v", true, retVal)
	}
}
