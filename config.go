package getconfig

import (
	"io/ioutil"
	"log"
)

/*Config returns config key value map to be used by program*/
func Config(configFilePath string) map[string]interface{} {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	noCommentsData := RemoveComments(data)
	congigMap := Map(configFilePath, noCommentsData)

	return congigMap
}
