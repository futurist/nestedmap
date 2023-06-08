# nestedmap

[![PkgGoDev](https://pkg.go.dev/badge/github.com/futurist/nestedmap)](https://pkg.go.dev/github.com/futurist/nestedmap)
[![Build Status](https://github.com/futurist/nestedmap/workflows/CI/badge.svg)](https://github.com/futurist/nestedmap/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/futurist/nestedmap)](https://goreportcard.com/report/github.com/futurist/nestedmap)
![Coverage](https://github.com/futurist/nestedmap/blob/main/.github/badge.svg?branch=badge)

golang data structure for manipulating and representing JSON objects with deeply nested map support.

More details please check [godoc](https://pkg.go.dev/github.com/futurist/nestedmap)

## Install

```sh
go get github.com/futurist/nestedmap
```

## Usage

```go
// below is true:
import "github.com/futurist/nestedmap"

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
_ := json.Unmarshal([]byte(input), &nestedMap)

nestedMap.GetValue("[A][B][C][D]")
// Output: value

nestedMap.SetValue("[A][B][C][E]", "OK") // OK
```
