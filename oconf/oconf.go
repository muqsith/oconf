package oconf

import (
	"encoding/json"
	"log"
	"path/filepath"
)

var filesMap map[string]bool

func setFileLoadedTrue(filePath string) {
	filesMap[filePath] = true
}

func isFileLoaded(filePath string) bool {
	result, ok := filesMap[filePath]
	if !ok {
		result = false
	}
	return result
}

func getAllConfig(configFilePath string) map[string]interface{} {
	filesMap = make(map[string]bool)
	setFileLoadedTrue(configFilePath)
	noCommentsData := getFileDataWithoutComments(configFilePath)
	dirPrefix := filepath.Dir(configFilePath)
	configMap := getConfigMap(dirPrefix, noCommentsData)
	return configMap
}

/*GetConfig returns a map of config*/
func GetConfig(configFilePath string) map[string]interface{} {
	configMap := getAllConfig(configFilePath)
	removePublicAnnotations(configMap)
	return configMap
}

/*GetFlatConfig returns a map of config with nested values separated by dot eg: pin := map["country.city.pin"]*/
func GetFlatConfig(configFilePath string) map[string]interface{} {
	configMap := getAllConfig(configFilePath)
	removePublicAnnotations(configMap)
	flatMap := make(map[string]interface{})
	createFlatMap("", flatMap, configMap)
	return flatMap
}

/*GetConfigAsJSONString returns JSON string of config*/
func GetConfigAsJSONString(configFilePath string) string {
	configMap := getAllConfig(configFilePath)
	removePublicAnnotations(configMap)
	jsonBytes, jsonErr := json.Marshal(configMap)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return string(jsonBytes)
}

/*GetPublicConfig returns a map of config*/
func GetPublicConfig(configFilePath string) map[string]interface{} {
	result := getAllConfig(configFilePath)
	saveOnlyWithPublicAnnotations(result)
	removePublicAnnotations(result)
	return result
}

/*GetPublicConfigAsJSONString returns a map of config*/
func GetPublicConfigAsJSONString(configFilePath string) string {
	result := getAllConfig(configFilePath)
	saveOnlyWithPublicAnnotations(result)
	removePublicAnnotations(result)
	jsonBytes, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return string(jsonBytes)
}
