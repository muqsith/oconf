package oconf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// remember to add -t 0 flag  like below
// go test -timeout 0 -run TestGetConfig
// go test -run TestGetConfig
func TestGetConfig(t *testing.T) {
	if configMap := GetConfig("./test-data/specific.cjson"); len(configMap) == 0 {
		t.Error("Not your config")
	} else {
		fmt.Printf("Map size: %d\n", len(configMap))
		for k, v := range configMap {
			fmt.Printf("Key: %s, Value: %v\n", k, v)
		}
		assert.Equal(t, "Walldorf", configMap["city"])
	}
}

// go test -run TestGetConfigAsJSONString
func TestGetConfigAsJSONString(t *testing.T) {
	configJsonString := GetConfigAsJSONString("./test-data/specific.cjson")
	fmt.Println(configJsonString)
	assert.Equal(t, true, len(configJsonString) > 0)
}

// go test -run TestGetMergedMap
func TestGetMergedMap(t *testing.T) {

	map1 := make(map[string]interface{})
	map2 := make(map[string]interface{})

	map1["address"] = map[string]interface{}{
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
	}

	map2["address"] = map[string]interface{}{
		"Street": "Street 5",
		"Area":   "Al Qusais",
		"City": map[string]interface{}{
			"province": "Al Qusais - Industrial Area",
		},
		"State": "Dubai",
		"Country": map[string]interface{}{
			"Name": "United Arab Emirates",
		},
	}

	merged := getMergedMap(map1, map2)
	address := merged["address"].(map[string]interface{})
	country := address["Country"].(map[string]interface{})
	countryCode := country["CountryCode"].(int)
	city := address["City"].(map[string]interface{})
	province := city["province"].(string)
	assert.Equal(t, 971, countryCode)
	assert.Equal(t, "Al Qusais - Industrial Area", province)
}

// go test -run TestGetPublicConfig
func TestGetPublicConfig(t *testing.T) {
	publicConfigMap := GetPublicConfig("./test-data/specific.cjson")
	fmt.Printf("Map size: %d\n", len(publicConfigMap))
	assert.Equal(t, "Baden-Wuttenberg", publicConfigMap["state"])
}

// go test -run TestGetPublicConfigAsJSONString
func TestGetPublicConfigAsJSONString(t *testing.T) {
	publicConfigString := GetPublicConfigAsJSONString("./test-data/specific.cjson")
	fmt.Println(publicConfigString)
	assert.Equal(t, true, len(publicConfigString) > 0)
}

// go test -run TestGetFlatConfig
func TestGetFlatConfig(t *testing.T) {
	publicConfigString := GetFlatConfig("./test-data/specific.cjson")
	fmt.Println(publicConfigString)
	assert.Equal(t, true, len(publicConfigString) > 0)
}
