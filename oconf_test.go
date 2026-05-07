package oconf

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetConfig1 tests basic configuration loading
// Run with: go test -run TestGetConfig1
func TestGetConfig1(t *testing.T) {
	configMap, err := GetConfig("test-data/specific.cjson")
	require.NoError(t, err, "GetConfig should not return an error")
	require.NotEmpty(t, configMap, "Config map should not be empty")

	fmt.Printf("Map size: %d\n", len(configMap))
	for k, v := range configMap {
		fmt.Printf("Key: %s, Value: %v\n", k, v)
	}

	assert.Equal(t, "Walldorf", configMap["city"])
}

// TestGetConfigAsJSONString tests JSON string output
// Run with: go test -run TestGetConfigAsJSONString
func TestGetConfigAsJSONString(t *testing.T) {
	configJsonString, err := GetConfigAsJSONString("./test-data/specific.cjson")
	require.NoError(t, err, "GetConfigAsJSONString should not return an error")

	fmt.Println(configJsonString)
	assert.Greater(t, len(configJsonString), 0, "JSON string should not be empty")
}

// TestGetMergedMap tests the map merging functionality
// Run with: go test -run TestGetMergedMap
func TestGetMergedMap(t *testing.T) {
	map1 := map[string]interface{}{
		"address": map[string]interface{}{
			"Street": "Street 2",
			"Area":   "Discovery Gardens",
			"City": map[string]interface{}{
				"name":     "Dubai",
				"province": "Jabal Ali",
			},
			"State": "Dubai",
			"Country": map[string]interface{}{
				"CountryCode": 971,
				"Name":        "United Arab Emirates",
			},
		},
	}

	map2 := map[string]interface{}{
		"address": map[string]interface{}{
			"Street": "Street 5",
			"Area":   "Al Qusais",
			"City": map[string]interface{}{
				"province": "Al Qusais - Industrial Area",
			},
			"State": "Dubai",
			"Country": map[string]interface{}{
				"Name": "United Arab Emirates",
			},
		},
	}

	merged := getMergedMap(map1, map2)
	address := merged["address"].(map[string]interface{})
	country := address["Country"].(map[string]interface{})

	// CountryCode should be preserved from map1
	countryCode, ok := country["CountryCode"].(int)
	require.True(t, ok, "CountryCode should be an int")
	assert.Equal(t, 971, countryCode)

	// Province should be overridden by map2
	city := address["City"].(map[string]interface{})
	province := city["province"].(string)
	assert.Equal(t, "Al Qusais - Industrial Area", province)
}

// TestGetPublicConfig tests filtering of public-only configuration
// Run with: go test -run TestGetPublicConfig
func TestGetPublicConfig(t *testing.T) {
	publicConfigMap, err := GetPublicConfig("./test-data/specific.cjson")
	require.NoError(t, err, "GetPublicConfig should not return an error")

	fmt.Printf("Map size: %d\n", len(publicConfigMap))
	assert.Equal(t, "Baden-Wuttenberg", publicConfigMap["state"])

	// Verify that arrays inside a #public parent are included (regression test)
	bag, ok := publicConfigMap["bag"].(map[string]interface{})
	require.True(t, ok, "bag should be present in public config")
	veggies, ok := bag["veggies"].(map[string]interface{})
	require.True(t, ok, "veggies should be present under bag in public config")
	leafy, ok := veggies["leafy"].([]interface{})
	require.True(t, ok, "leafy array should be present inside veggies#public")
	assert.Equal(t, []interface{}{"spinach", "basil"}, leafy)
}

// TestGetPublicConfigAsJSONString tests public config as JSON string
// Run with: go test -run TestGetPublicConfigAsJSONString
func TestGetPublicConfigAsJSONString(t *testing.T) {
	publicConfigString, err := GetPublicConfigAsJSONString("./test-data/specific.cjson")
	require.NoError(t, err, "GetPublicConfigAsJSONString should not return an error")

	fmt.Println(publicConfigString)
	assert.Greater(t, len(publicConfigString), 0, "Public config JSON should not be empty")
}

// TestGetFlatConfig tests flat configuration format
// Run with: go test -run TestGetFlatConfig
func TestGetFlatConfig(t *testing.T) {
	flatConfig, err := GetFlatConfig("./test-data/specific.cjson")
	require.NoError(t, err, "GetFlatConfig should not return an error")

	fmt.Println(flatConfig)
	assert.Greater(t, len(flatConfig), 0, "Flat config should not be empty")

	// Verify flat format works (nested keys joined with dots)
	// The actual keys depend on the test data structure
	for key := range flatConfig {
		fmt.Printf("Flat key: %s\n", key)
	}
}

// TestGetFlatPublicConfig tests flat public configuration format
// Run with: go test -run TestGetFlatPublicConfig
func TestGetFlatPublicConfig(t *testing.T) {
	flatPublicConfig, err := GetFlatPublicConfig("./test-data/specific.cjson")
	require.NoError(t, err, "GetFlatPublicConfig should not return an error")

	fmt.Println(flatPublicConfig)
	assert.Greater(t, len(flatPublicConfig), 0, "Flat public config should not be empty")
}

// TestGetConfigSubDir tests configuration loading from subdirectories with includes
// Run with: go test -run TestGetConfigSubDir
func TestGetConfigSubDir(t *testing.T) {
	configMap, err := GetConfig("./test-data/sub1-dir/sub2-dir2/ref1.cjson")
	require.NoError(t, err, "GetConfig should not return an error")
	require.NotEmpty(t, configMap, "Config map should not be empty")

	fmt.Printf("Map size: %d\n", len(configMap))
	for k, v := range configMap {
		fmt.Printf("Key: %s, Value: %v\n", k, v)
	}

	assert.Equal(t, "Hidelberg", configMap["city"])
}

// TestBadConfig tests error handling for invalid configuration
// Run with: go test -run TestBadConfig
func TestBadConfig(t *testing.T) {
	_, err := GetFlatConfig("./test-data/bad.cjson")
	assert.Error(t, err, "GetFlatConfig should return an error for bad config")
	fmt.Printf("Expected error: %v\n", err)
}

// TestFileNotFound tests error handling for missing files
// Run with: go test -run TestFileNotFound
func TestFileNotFound(t *testing.T) {
	_, err := GetConfig("./test-data/nonexistent.cjson")
	assert.Error(t, err, "GetConfig should return an error for missing file")
	fmt.Printf("Expected error: %v\n", err)
}

// TestInvalidJSON tests error handling for invalid JSON
// Run with: go test -run TestInvalidJSON
func TestInvalidJSON(t *testing.T) {
	// This test requires a test file with invalid JSON
	// Create test-data/invalid.cjson with invalid JSON if needed
	_, err := GetConfig("./test-data/invalid.cjson")
	if err != nil {
		fmt.Printf("Expected error for invalid JSON: %v\n", err)
	}
	// Note: This test will only work if invalid.cjson exists
}

// TestCreateFlatMap tests the flat map creation with various types
// Run with: go test -run TestCreateFlatMap
func TestCreateFlatMap(t *testing.T) {
	nested := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"value": "deep",
			},
			"number": 42.0,
		},
		"array":  []interface{}{"a", "b", "c"},
		"string": "test",
		"bool":   true,
	}

	flat := make(map[string]interface{})
	err := createFlatMap("", flat, nested)
	require.NoError(t, err, "createFlatMap should not return an error")

	assert.Equal(t, "deep", flat["level1.level2.value"])
	assert.Equal(t, 42.0, flat["level1.number"])
	assert.Equal(t, "test", flat["string"])
	assert.Equal(t, true, flat["bool"])
	assert.NotNil(t, flat["array"])

	fmt.Println("Flattened map:")
	for k, v := range flat {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

// TestIsMap tests the type checking function
// Run with: go test -run TestIsMap
func TestIsMap(t *testing.T) {
	assert.True(t, isMap(map[string]interface{}{}))
	assert.True(t, isMap(map[string]interface{}{"key": "value"}))
	assert.False(t, isMap("string"))
	assert.False(t, isMap(42))
	assert.False(t, isMap(nil))
	assert.False(t, isMap([]interface{}{}))
}

// TestNestedIncludeSingle tests a single #include inside a nested object.
// Keys in the nested object override keys from the included file.
// Run with: go test -run TestNestedIncludeSingle
func TestNestedIncludeSingle(t *testing.T) {
	configMap, err := GetConfig("./test-data/nested-include/single-nested.cjson")
	require.NoError(t, err)

	assert.Equal(t, "myapp", configMap["app"])

	db, ok := configMap["database"].(map[string]interface{})
	require.True(t, ok, "database should be a map")
	assert.Equal(t, "db.example.com", db["host"])
	assert.Equal(t, float64(5432), db["port"])
	// "name" in the nested object overrides the included file's "name"
	assert.Equal(t, "override-db", db["name"])
}

// TestNestedIncludeMultipleFiles tests merging multiple files inside a nested #include.
// Files are merged left-to-right; the last file wins on key conflicts.
// Run with: go test -run TestNestedIncludeMultipleFiles
func TestNestedIncludeMultipleFiles(t *testing.T) {
	configMap, err := GetConfig("./test-data/nested-include/multi-files-nested.cjson")
	require.NoError(t, err)

	assert.Equal(t, "myapp", configMap["app"])

	services, ok := configMap["services"].(map[string]interface{})
	require.True(t, ok, "services should be a map")

	// from db.cjson only (cache.cjson has no "name")
	assert.Equal(t, "mydb", services["name"])
	// from cache.cjson only (db.cjson has no "ttl" or "maxConn")
	assert.Equal(t, float64(300), services["ttl"])
	assert.Equal(t, float64(100), services["maxConn"])
	// "host" exists in both; cache.cjson (second) wins
	assert.Equal(t, "cache.example.com", services["host"])
}

// TestNestedIncludeSiblings tests that multiple sibling nested objects can each have their
// own independent #include directives that are resolved correctly.
// Run with: go test -run TestNestedIncludeSiblings
func TestNestedIncludeSiblings(t *testing.T) {
	configMap, err := GetConfig("./test-data/nested-include/sibling-nested.cjson")
	require.NoError(t, err)

	assert.Equal(t, "myapp", configMap["app"])

	db, ok := configMap["database"].(map[string]interface{})
	require.True(t, ok, "database should be a map")
	assert.Equal(t, "db.example.com", db["host"])
	assert.Equal(t, float64(5432), db["port"])
	assert.Equal(t, "mydb", db["name"])

	server, ok := configMap["server"].(map[string]interface{})
	require.True(t, ok, "server should be a map")
	assert.Equal(t, "api.example.com", server["host"])
	assert.Equal(t, true, server["tls"])
	// local "port" overrides the included server.cjson port of 8080
	assert.Equal(t, float64(9000), server["port"])
}

// TestNestedIncludeDeep tests #include inside a doubly-nested object.
// Run with: go test -run TestNestedIncludeDeep
func TestNestedIncludeDeep(t *testing.T) {
	configMap, err := GetConfig("./test-data/nested-include/deep-nested.cjson")
	require.NoError(t, err)

	assert.Equal(t, "myapp", configMap["app"])

	services, ok := configMap["services"].(map[string]interface{})
	require.True(t, ok, "services should be a map")

	backend, ok := services["backend"].(map[string]interface{})
	require.True(t, ok, "services.backend should be a map")
	assert.Equal(t, "db.example.com", backend["host"])
	assert.Equal(t, float64(5432), backend["port"])
	assert.Equal(t, "mydb", backend["name"])
}

// TestNestedIncludeWithRootInclude tests that root-level and nested #include directives
// in the same file both work independently and do not interfere.
// Run with: go test -run TestNestedIncludeWithRootInclude
func TestNestedIncludeWithRootInclude(t *testing.T) {
	configMap, err := GetConfig("./test-data/nested-include/with-root-include.cjson")
	require.NoError(t, err)

	// root-level keys from server.cjson
	assert.Equal(t, "api.example.com", configMap["host"])
	assert.Equal(t, float64(8080), configMap["port"])
	assert.Equal(t, true, configMap["tls"])

	// nested include from db.cjson
	db, ok := configMap["database"].(map[string]interface{})
	require.True(t, ok, "database should be a map")
	assert.Equal(t, "db.example.com", db["host"])
	assert.Equal(t, float64(5432), db["port"])
	assert.Equal(t, "mydb", db["name"])
}

// TestNestedIncludeChained tests that an included file which itself uses root-level
// #include has all of its own includes resolved before being merged into the parent.
// Run with: go test -run TestNestedIncludeChained
func TestNestedIncludeChained(t *testing.T) {
	// chained-nested.cjson → database.#include → db-with-defaults.cjson
	//                                              which has root #include → conn-defaults.cjson
	configMap, err := GetConfig("./test-data/nested-include/chained-nested.cjson")
	require.NoError(t, err)

	assert.Equal(t, "myapp", configMap["app"])

	db, ok := configMap["database"].(map[string]interface{})
	require.True(t, ok, "database should be a map")

	// from db-with-defaults.cjson directly
	assert.Equal(t, "db.example.com", db["host"])
	assert.Equal(t, float64(5432), db["port"])

	// from conn-defaults.cjson via db-with-defaults.cjson's root #include
	assert.Equal(t, float64(30), db["timeout"])
	assert.Equal(t, float64(3), db["retries"])
}

// TestNestedIncludeErrorBadType tests that a non-array #include value inside a nested
// object produces an immediate TypeError and does not silently skip.
// Run with: go test -run TestNestedIncludeErrorBadType
func TestNestedIncludeErrorBadType(t *testing.T) {
	_, err := GetConfig("./test-data/nested-include/bad-include-type.cjson")
	require.Error(t, err)

	var typeErr *TypeError
	require.True(t, errors.As(err, &typeErr), "expected TypeError, got: %v", err)
	assert.Equal(t, IncludeKey, typeErr.Key)
	assert.Equal(t, "array", typeErr.Expected)

	// error message must include the key path so the caller knows where it failed
	assert.True(t, strings.Contains(err.Error(), "database"), "error should name the failing key: %v", err)
}

// TestNestedIncludeErrorMissingFile tests that referencing a nonexistent file inside a
// nested #include fails immediately with a FileError containing the bad path.
// Run with: go test -run TestNestedIncludeErrorMissingFile
func TestNestedIncludeErrorMissingFile(t *testing.T) {
	_, err := GetConfig("./test-data/nested-include/missing-file-ref.cjson")
	require.Error(t, err)

	var fileErr *FileError
	require.True(t, errors.As(err, &fileErr), "expected FileError, got: %v", err)
	assert.Contains(t, fileErr.Path, "nonexistent.cjson")

	// error message must include the key path so the caller knows where it failed
	assert.True(t, strings.Contains(err.Error(), "database"), "error should name the failing key: %v", err)
}
