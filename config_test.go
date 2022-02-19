package getconfig

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// remember to add -t 0 flag  like below
// go test -timeout 0 -run TestConfig
// go test -run TestConfig
func TestConfig(t *testing.T) {
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
