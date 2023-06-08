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

func ExampleInReadme() {
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
