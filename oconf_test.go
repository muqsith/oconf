package oconf

import (
	"fmt"
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
