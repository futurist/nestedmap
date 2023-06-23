package nestedmap_test

import (
	"encoding/json"
	"testing"

	nestedmap "github.com/futurist/nestedmap"
)

func TestGetValue(t *testing.T) {
	t.Parallel()

	nMap := &nestedmap.NestedMap{
		Data: make(map[string]interface{}),
	}

	success := nMap.SetValue("[A][B][C][E][X]", "xyz")
	val := nMap.GetValue("[A][B][C][E][X]")

	if success != true {
		t.Errorf("Value set failed")
	}
	if val != "xyz" {
		t.Errorf("Expected value: xyz, found: %v", val)
	}
}

func TestCustomSerialization(t *testing.T) {
	t.Parallel()

	nMap := &nestedmap.NestedMap{
		Data: make(map[string]interface{}),
	}

	nMap.SetValue("[A][B][C]", "value")
	serialized, err := json.Marshal(nMap)

	if err != nil {
		t.Errorf("Error marhsaling: %v", err)
	}
	expectedOutput := `{"A":{"B":{"C":"value"}}}`
	if string(serialized) != expectedOutput {
		t.Errorf("Expected: %v, found: %v", expectedOutput, string(serialized))
	}
}
