package oconf

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getFromIncludes(dirPrefix string, includedFiles []interface{}) map[string]interface{} {
	pathSep := string(os.PathSeparator)
	mNew := make(map[string]interface{})
	for _, includedFile := range includedFiles {
		includedFilePath := includedFile.(string)
		if filepath.IsAbs(includedFilePath) {
			dirPrefix = filepath.Dir(includedFilePath)
		} else {
			if strings.Contains(includedFilePath, pathSep) {
				includedFilePath = path.Join(dirPrefix, includedFilePath)
				dirPrefix = filepath.Dir(includedFilePath)
			} else {
				includedFilePath = dirPrefix + pathSep + includedFilePath
			}
		}
		if !isFileLoaded(includedFilePath) {
			noCommentsData := getFileDataWithoutComments(includedFilePath)
			mapFromInclude := getConfigMap(dirPrefix, noCommentsData)
			if len(mNew) > 1 {
				mNew = getMergedMap(mNew, mapFromInclude)
			} else {
				mNew = mapFromInclude
			}
		}
	}
	return mNew
}

func getFileDataWithoutComments(configFilePath string) []byte {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	noCommentsData := removeComments(data)
	return noCommentsData
}

/*
getConfigMap returns the map. If `returnFlatMap` is true then it returns map with nested keys separated by dots.
Eg: dev.host.ip : 10.127.1.4
*/
func getConfigMap(dirPrefix string, data []byte) map[string]interface{} {
	var result map[string]interface{}
	var f interface{}
	jsonErr := json.Unmarshal(data, &f)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	result = f.(map[string]interface{})
	includes, ok := result["#include"]
	if ok {
		includedFiles := includes.([]interface{})
		includesConfig := getFromIncludes(dirPrefix, includedFiles)
		result = getMergedMap(includesConfig, result)
		delete(result, "#include")
	}
	return result
}
