package oconf

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func isMap(v interface{}) bool {
	result := v != nil
	t := fmt.Sprintf("%T", v)
	result = result && strings.Contains(t, "map[string]")
	return result
}

func removePublicAnnotations(configMap map[string]interface{}) {
	for key, val := range configMap {
		if strings.Contains(key, "#public") {
			newKey := strings.Replace(key, "#public", "", -1)
			delete(configMap, key)
			configMap[newKey] = val
		}
		if isMap(val) {
			removePublicAnnotations(val.(map[string]interface{}))
		}
	}
}

func saveOnlyWithPublicAnnotations(configMap map[string]interface{}) {
	for key, val := range configMap {
		if !isMap(val) {
			if !strings.Contains(key, "#public") {
				delete(configMap, key)
			}
		}

		if isMap(val) && !strings.Contains(key, "#public") {
			subMap := val.(map[string]interface{})
			saveOnlyWithPublicAnnotations(subMap)
			if len(subMap) == 0 {
				delete(configMap, key)
			}
		}

	}
}

func createFlatMap(prefix string, newMap, configMap map[string]interface{}) {
	for k, v := range configMap {
		keyPrefix := prefix
		if len(prefix) > 0 {
			keyPrefix = fmt.Sprintf("%s%s", prefix, ".")
		}
		key := fmt.Sprintf("%s%s", keyPrefix, k)

		switch vv := v.(type) {
		case string, float64, bool, []interface{}:
			{
				newMap[key] = v
			}
		case interface{}:
			{
				createFlatMap(key, newMap, v.(map[string]interface{}))
			}
		default:
			log.Fatalf("Key:%s has value of a %T type I don't know how to handle", k, vv)
		}
	}
}

func parseDate(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func calculateDaysBetweenDates(startDate, endDate string) int {
	startDateTime := parseDate(startDate)
	endDateTime := parseDate(endDate)
	return int(endDateTime.Sub(startDateTime).Hours() / 24)
}
