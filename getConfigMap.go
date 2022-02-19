package getconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func copyConfigMap(m1, m2 map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(m1)-1+len(m2))
	for k1, v1 := range m1 {
		m[k1] = v1
	}
	for k1, v1 := range m2 {
		m[k1] = v1
	}
	return m
}

func getFlatMap(configFilePath, prefix string, m map[string]interface{}) map[string]interface{} {

	pathSep := string(os.PathSeparator)

	mapSize := len(m)
	mNew := make(map[string]interface{}, mapSize)

	includes, ok := m["#include"]

	if ok {
		includedFiles := includes.([]interface{})
		for _, includedFile := range includedFiles {
			includedFilePath := includedFile.(string)
			if !filepath.IsAbs(includedFilePath) {
				includedFilePath = filepath.Dir(configFilePath) + pathSep + includedFilePath
			}
			mapFromInclude := GetConfig(includedFilePath)
			mNew = copyConfigMap(mNew, mapFromInclude)
		}
		delete(m, "#include")
	}

	for k, v := range m {
		switch vv := v.(type) {
		case string, float64, bool, []interface{}:
			{
				keyPrefix := prefix
				if len(prefix) > 0 {
					keyPrefix = fmt.Sprintf("%s%s", prefix, ".")
				}
				key := fmt.Sprintf("%s%s", keyPrefix, k)
				mNew[key] = v
			}
		case interface{}:
			{
				mNested := getFlatMap(configFilePath, k, v.(map[string]interface{}))
				mNew = copyConfigMap(mNew, mNested)
			}
		default:
			log.Fatalf("Key:%s has value of a %T type I don't know how to handle", k, vv)
		}
	}

	return mNew
}

/*
ConfigMap returns the map of the final configuration. Keys are separated by dots
Eg: dev.host.ip : 10.127.1.4
*/
func ConfigMap(configFilePath string, data []byte) map[string]interface{} {
	var f interface{}
	jsonErr := json.Unmarshal(data, &f)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	m := f.(map[string]interface{})

	flatMap := getFlatMap(configFilePath, "", m)

	return flatMap
}
