package treeutil

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func DeepEqualNormalized(a interface{}, b interface{}) bool {
	normalize := func(v interface{}) interface{} {
		bytes, err := json.Marshal(v)
		if err != nil {
			return v
		}

		var result interface{}
		if err := json.Unmarshal(bytes, &result); err != nil {
			return v
		}

		return result
	}

	return reflect.DeepEqual(normalize(a), normalize(b))
}

func GetValueAtPath(obj map[string]interface{}, path string) (interface{}, bool) {
	var current interface{} = obj
	decodeEscapeSequences := func(part string) string {
		// Replace escape sequences according to RFC 6901
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")
		return part
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := range parts {
		parts[i] = decodeEscapeSequences(parts[i])
	}

	for _, part := range parts {
		switch curr := current.(type) {

		case map[string]interface{}:
			val, ok := curr[part]
			if !ok {
				return nil, false
			}
			current = val

		case []interface{}:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(curr) {
				return nil, false
			}
			current = curr[idx]

		default:
			return nil, false
		}
	}

	return current, true
}
