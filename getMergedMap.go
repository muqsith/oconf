package oconf

/*
	getMergedMap map2 overrides map1
*/
func getMergedMap(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key1, value1 := range map1 {
		result[key1] = value1
	}
	for key2, value2 := range map2 {
		// copy value2 if it doesn't exist already
		if _, ok := result[key2]; !ok {
			result[key2] = value2
		} else {
			// override if not map
			if isMap(value2) {
				v1 := result[key2] // value from map1
				result[key2] = getMergedMap(v1.(map[string]interface{}), value2.(map[string]interface{}))
			} else {
				result[key2] = value2
			}
		}

	}
	return result
}