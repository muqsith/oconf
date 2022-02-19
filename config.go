package getconfig

import (
	"io/ioutil"
	"log"
)

/*GetConfig returns a map of config, use . for nested properties eg: address.country.city*/
func GetConfig(configFilePath string) map[string]interface{} {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	noCommentsData := RemoveComments(data)
	congigMap := ConfigMap(configFilePath, noCommentsData)

	return congigMap
}
