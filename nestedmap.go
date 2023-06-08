// Package nestedmap implements a data structure for manipulating and representing JSON objects
// with a depth of nesting. The NestedMap type provides methods for setting and getting values
// by specifying a key-path string that describes the location of the key-value pair within the JSON object.
package nestedmap

import (
	"encoding/json"
	"regexp"
	"strconv"
)

// NestedMap is a data structure that represents a map with nested keys as a tree structure.
type NestedMap struct {
	Data map[string]interface{}
}

func New() *NestedMap {
	return &NestedMap{
		Data: make(map[string]interface{}),
	}
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

// MarshalJSON is a custom marshaller for NestedMap that
// serializes the NestedMap without exposing the `Data` structure.
func (n *NestedMap) MarshalJSON() ([]byte, error) {
	processedMap := n.processMap(n.Data)
	return json.Marshal(processedMap)
}

// processMap recursively converts NestedMap values to regular maps,
// excluding the `Data` structure, and returns the result as a map[string]interface{}.
func (n *NestedMap) processMap(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(m))
	for key, value := range m {
		switch v := value.(type) {
		case *NestedMap:
			res[key] = n.processMap(v.Data)
		case []*NestedMap:
			res[key] = n.processSlice(v)
		default:
			res[key] = value
		}
	}
	return res
}

// processSlice processes a slice of NestedMap to handle serialization.
func (n *NestedMap) processSlice(nestedMaps []*NestedMap) []interface{} {
	res := make([]interface{}, len(nestedMaps))
	for i, nm := range nestedMaps {
		res[i] = n.processMap(nm.Data)
	}
	return res
}

// GetValue returns the value at the specified path in the NestedMap.
// If the value is not found at the specified path, it returns nil.
func (n *NestedMap) GetValue(path string) interface{} {
	keys := parsePath(path)
	if keys == nil {
		return nil
	}

	value, found := n.getValueHelper(keys, 0)
	if !found {
		return nil
	}

	return value
}

// SetValue sets the value at the specified path in the NestedMap.
// It will create new NestedMaps along the path if they don't exist.
// The function returns true if the value is set successfully, false otherwise.
func (n *NestedMap) SetValue(path string, value interface{}) bool {
	keys := parsePath(path)
	if keys == nil {
		return false
	}

	return n.setValueHelper(keys, 0, value)
}

// getValueHelper is a helper function to search for the value at the specified path.
func (n *NestedMap) getValueHelper(keys []string, index int) (interface{}, bool) {
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

	switch v := value.(type) {
	case *NestedMap:
		return v.getValueHelper(keys, index+1)
	case []*NestedMap:
		indexInSlice, err := strconv.Atoi(keys[index+1])
		if err != nil || indexInSlice >= len(v) {
			return nil, false
		}
		return v[indexInSlice].getValueHelper(keys, index+2)
	default:
		return nil, false
	}
}

// setValueHelper is a helper function to set the value at the specified path.
// It returns true if the value is set successfully, false otherwise.
func (n *NestedMap) setValueHelper(keys []string, index int, value interface{}) bool {
	if index >= len(keys) {
		return false
	}

	if index == len(keys)-1 {
		n.Data[keys[index]] = value
		return true
	}

	if currentValue, ok := n.Data[keys[index]].(*NestedMap); ok {
		return currentValue.setValueHelper(keys, index+1, value)
	}

	newMap := New()
	n.Data[keys[index]] = newMap
	return newMap.setValueHelper(keys, index+1, value)
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
