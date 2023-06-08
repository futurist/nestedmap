// Package nestedmap implements a data structure for manipulating and representing JSON objects
// with a depth of nesting. The NestedMap type provides methods for setting and getting values
// by specifying a key-path string that describes the location of the key-value pair within the JSON object.
package nestedmap

import (
	"encoding/json"
	"regexp"
)

// NestedMap is a data structure that represents a map with nested keys as a tree structure.
type NestedMap struct {
	Data map[string]interface{}
}

// UnmarshalJSON is a custom unmarshaller for NestedMap that translates the JSON object
// into a nested map structure and fixes the nested maps to be of type NestedMap.
func (n *NestedMap) UnmarshalJSON(data []byte) error {
	tmpMap := make(map[string]interface{})
	err := json.Unmarshal(data, &tmpMap)
	if err != nil {
		return err
	}

	n.fixNestedMaps(tmpMap)

	n.Data = tmpMap
	return nil
}

// fixNestedMaps recursively converts nested maps into NestedMap types.
func (n *NestedMap) fixNestedMaps(m map[string]interface{}) {
	for key, value := range m {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			tmpNestedMap := &NestedMap{}
			tmpNestedMap.fixNestedMaps(nestedMap)
			tmpNestedMap.Data = nestedMap
			m[key] = tmpNestedMap
		}
	}
}

// GetValue returns the value at the specified path in the NestedMap.
// If the value is not found at the specified path, it returns nil.
func (n NestedMap) GetValue(path string) interface{} {
	keys := parsePath(path)
	if keys == nil {
		return nil
	}

	value, found := n.getValue_helper(keys, 0)
	if !found {
		return nil
	}

	return value
}

// SetValue sets the value at the specified path in the NestedMap.
// It will create new NestedMaps along the path if they don't exist.
func (n NestedMap) SetValue(path string, value interface{}) {
	keys := parsePath(path)
	if keys == nil {
		return
	}

	n.setValue_helper(keys, 0, value)
}

// getValue_helper is a helper function to search for the value at the specified path.
func (n NestedMap) getValue_helper(keys []string, index int) (interface{}, bool) {
	if index >= len(keys) {
		return nil, false
	}

	value, exists := n.Data[keys[index]]
	if !exists {
		return nil, false
	}

	if index == len(keys)-1 {
		return value, true
	}

	if nestedMap, ok := value.(*NestedMap); ok {
		return nestedMap.getValue_helper(keys, index+1)
	}

	return nil, false
}

// setValue_helper is a helper function to set the value at the specified path.
func (n NestedMap) setValue_helper(keys []string, index int, value interface{}) {
	if index >= len(keys) {
		return
	}

	if index == len(keys)-1 {
		n.Data[keys[index]] = value
		return
	}

	if currentValue, ok := n.Data[keys[index]].(*NestedMap); ok {
		currentValue.setValue_helper(keys, index+1, value)
	} else {
		newMap := &NestedMap{
			Data: make(map[string]interface{}),
		}
		n.Data[keys[index]] = newMap
		newMap.setValue_helper(keys, index+1, value)
	}
}

var pathRE = regexp.MustCompile(`\[(.+?)\]`)

// parsePath returns an array of keys extracted from the given path string.
func parsePath(path string) []string {
	matches := pathRE.FindAllStringSubmatch(path, -1)

	if matches == nil {
		return nil
	}

	keys := make([]string, 0, len(matches))
	for _, match := range matches {
		keys = append(keys, match[1])
	}

	return keys
}
