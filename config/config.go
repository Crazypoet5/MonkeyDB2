package config

import (
    "../common"
    "encoding/json"
    "io/ioutil"
)

func LoadConfig(c string) map[string]interface{} {
    b, err := ioutil.ReadFile(common.FixPath(common.GetCurrentDirectory() + "/etc/" + c + ".json"))
    if err != nil {
        return nil
    }
    var m map[string]interface{}
    json.Unmarshal(b, &m)
    return m
}