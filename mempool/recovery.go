package mempool

import (
    "encoding/json"
    "io/ioutil"
    "../log"
    "strconv"
)

func syncTableToFile() {
    ready := make(map[string]fileImage)
    for k, v := range MallocTable {
        str := strconv.Itoa(int(k))
        ready[str] = v
    }
    data, err := json.Marshal(ready)
    if err != nil {
        log.WriteLog("sys", "JSON stringfy malloc table with an error : " + err.Error())
        return
    }
    ioutil.WriteFile(".\\malloc.json", data, 0666)
}