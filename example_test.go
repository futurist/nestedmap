package nestedmap

import (
	"encoding/json"
	"fmt"
)

func Example() {
	input := `
	{
		"A": {
			"B": {
				"C": {
					"D": "old_value"
				}
			}
		}
	}`

	var nestedMap NestedMap
	err := json.Unmarshal([]byte(input), &nestedMap)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	value := nestedMap.GetValue(`[A][B][C][D]`)
	if value != nil {
		fmt.Println("Value at path [A][B][C][D]:", value)
	} else {
		fmt.Println("Path not found")
	}

	value = nestedMap.GetValue(`[A][B][C][DD]`)
	fmt.Println("Value at [A][B][C][DD]:", value)

	newNestedMap := NestedMap{
		Data: map[string]interface{}{
			"X": "xyz",
			"Y": "abc",
		},
	}

	nestedMap.SetValue(`[A][B][C][E][XX]`, &newNestedMap)
	value = nestedMap.GetValue(`[A][B][C][E][XX][X]`)
	fmt.Println("Value at path [A][B][C][E][XX][X]:", value)

	serialized, err := json.Marshal(nestedMap.Data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Serialized JSON:", string(serialized))

	// Output:
	// Value at path [A][B][C][D]: old_value
	// Value at [A][B][C][DD]: <nil>
	// Value at path [A][B][C][E][XX][X]: xyz
	// Serialized JSON: {"A":{"B":{"C":{"D":"old_value","E":{"XX":{"X":"xyz","Y":"abc"}}}}}}
}

// ExampleMarshal the example in README.md
func ExampleMarshal() {
	input := `
	{
		"A": {
			"B": {
				"C": {
					"D": "value"
				}
			}
		}
	}`

	var nestedMap NestedMap
	_ = json.Unmarshal([]byte(input), &nestedMap)

	fmt.Println(nestedMap.GetValue("[A][B][C][D]"))

	fmt.Println(nestedMap.SetValue("[A][B][C][E]", "OK"))

	serialized, _ := json.Marshal(nestedMap.Data)
	fmt.Println(string(serialized))
	// Output:
	// value
	// true
	// {"A":{"B":{"C":{"D":"value","E":"OK"}}}}
}

func ExampleNew() {
	nestedMap := New()
	nestedMap1 := &NestedMap{
		Data: make(map[string]interface{}),
	}
	nestedMap2 := &NestedMap{
		Data: make(map[string]interface{}),
	}

	_ = nestedMap1.SetValue("[A][B]", "value1")
	_ = nestedMap2.SetValue("[X][Y]", "value2")

	// Create an array of NestedMap with nestedMap1 and nestedMap2.
	nestedMaps := []*NestedMap{nestedMap1, nestedMap2}
	_ = nestedMap.SetValue("[U]", nestedMaps) // sets the "U" key to the nestedMaps slice

	fmt.Printf("Value at [U][0][A][B]: %v\n", nestedMap.GetValue("[U][0][A][B]"))
	fmt.Printf("Value at [U][1][X][Y]: %v\n", nestedMap.GetValue("[U][1][X][Y]"))

	serialized, err := json.Marshal(nestedMap)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Serialized JSON:", string(serialized))
	// Output:
	// Value at [U][0][A][B]: value1
	// Value at [U][1][X][Y]: value2
	// Serialized JSON: {"U":[{"A":{"B":"value1"}},{"X":{"Y":"value2"}}]}
}
