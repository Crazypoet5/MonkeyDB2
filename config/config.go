package config

import (
	"encoding/json"
	"io/ioutil"

	"../log"

	"../common"
)

func LoadConfig(c string) map[string]interface{} {
	b, err := ioutil.ReadFile(common.FixPath(common.COMMON_DIR + "/" + c + ".json"))
	if err != nil {
		return nil
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil
	}
	log.WriteLog("sys", m)
	return m
}
