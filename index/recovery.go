package index

import (
    "encoding/json"
    "io/ioutil"
    "../log"
    "strconv"
    "./csbt"
)

type IndexContent struct {
    Kind        int
    Database    string
    Table       string
    Key         string
    Base        uintptr
}

func syncTableToFile() {
    ready := make(map[string]IndexContent)
    for k, v := range IndexTable {
        str := strconv.Itoa(int(k))
        ready[str] = IndexContent {
            Kind:       v.Kind,
            Database:   v.Database,
            Table:      v.Table,
            Key:        v.Key,
            Base:       v.I.(*csbt.DCSBT).BaseAddr,
        }
    }
    data, err := json.Marshal(ready)
    if err != nil {
        log.WriteLog("sys", "JSON stringfy malloc table with an error : " + err.Error())
        return
    }
    ioutil.WriteFile(".\\index.json", data, 0666)
}

func recovery() {
    
}