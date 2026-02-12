package oconf

// getMergedMap merges two maps where map2 values override map1 values.
// For nested maps, it recursively merges them rather than replacing entirely.
// This allows for partial overrides of nested configuration structures.
//
// Example:
//
//	map1 := map[string]interface{}{
//	    "server": map[string]interface{}{
//	        "host": "localhost",
//	        "port": 8080,
//	    },
//	}
//	map2 := map[string]interface{}{
//	    "server": map[string]interface{}{
//	        "port": 9000,
//	    },
//	}
//	result := getMergedMap(map1, map2)
//	// result["server"] = {"host": "localhost", "port": 9000}
func getMergedMap(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy all values from map1 first
	for key, value := range map1 {
		result[key] = value
	}

	// Merge or override with values from map2
	for key, value2 := range map2 {
		value1, exists := result[key]

		if !exists {
			// Key doesn't exist in map1, just copy it
			result[key] = value2
		} else if isMap(value2) && isMap(value1) {
			// Both values are maps, merge them recursively
			result[key] = getMergedMap(
				value1.(map[string]interface{}),
				value2.(map[string]interface{}),
			)
		} else {
			// Override with map2's value (handles primitives and type mismatches)
			result[key] = value2
		}
	}

	return result
}
