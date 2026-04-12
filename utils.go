package oconf

import (
	"fmt"
	"strings"
)

// isMap checks if the given value is a map[string]interface{}
func isMap(v interface{}) bool {
	if v == nil {
		return false
	}
	_, ok := v.(map[string]interface{})
	return ok
}

// removePublicAnnotations removes the PublicKey annotation from all keys in the map recursively
func removePublicAnnotations(configMap map[string]interface{}) {
	for key, val := range configMap {
		if strings.Contains(key, PublicKey) {
			newKey := strings.ReplaceAll(key, PublicKey, "")
			delete(configMap, key)
			configMap[newKey] = val
			key = newKey
		}
		if isMap(val) {
			removePublicAnnotations(val.(map[string]interface{}))
		}
	}
}

// saveOnlyWithPublicAnnotations removes all keys that don't have the PublicKey annotation
// It recursively processes nested maps and removes empty maps after filtering
func saveOnlyWithPublicAnnotations(configMap map[string]interface{}) {
	saveOnlyWithPublicAnnotationsHelper(configMap, false)
}

// saveOnlyWithPublicAnnotationsHelper is the recursive implementation that tracks whether
// we are already inside a public subtree (insidePublic). When insidePublic is true, all
// children are kept unconditionally since they belong to a parent marked as public.
func saveOnlyWithPublicAnnotationsHelper(configMap map[string]interface{}, insidePublic bool) {
	for key, val := range configMap {
		keyIsPublic := strings.Contains(key, PublicKey)

		if insidePublic {
			// Already inside a public subtree — keep everything, just recurse into submaps
			if isMap(val) {
				saveOnlyWithPublicAnnotationsHelper(val.(map[string]interface{}), true)
			}
		} else if !isMap(val) {
			// Non-map value outside a public subtree: keep only if key is marked public
			if !keyIsPublic {
				delete(configMap, key)
			}
		} else {
			subMap := val.(map[string]interface{})
			// If this key is public, all its children are implicitly public
			saveOnlyWithPublicAnnotationsHelper(subMap, keyIsPublic)

			// Remove the parent key if it's not public and its subtree is now empty
			if !keyIsPublic && len(subMap) == 0 {
				delete(configMap, key)
			}
		}
	}
}

// createFlatMap converts a nested map into a flat map where nested keys are joined with dots
// For example: {"a": {"b": "value"}} becomes {"a.b": "value"}
func createFlatMap(prefix string, flatMap, nestedMap map[string]interface{}) error {
	for key, value := range nestedMap {
		// Build the full key with dot notation
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string, float64, bool, nil, []interface{}:
			flatMap[fullKey] = value
		case map[string]interface{}:
			if err := createFlatMap(fullKey, flatMap, v); err != nil {
				return err
			}
		default:
			return &TypeError{
				Key:      fullKey,
				Expected: "string, number, bool, array, or object",
				Actual:   fmt.Sprintf("%T", v),
			}
		}
	}
	return nil
}
