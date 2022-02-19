package getconfig

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	if d := Config("/home/mui/one/webshop-stats-go/config/config.cjson"); len(d) == 0 {
		t.Error("Not your config")
	} else {
		fmt.Printf("Map size: %d\n", len(d))
		for k, v := range d {
			fmt.Printf("Key: %s, Value: %v\n", k, v)
		}
	}
}
