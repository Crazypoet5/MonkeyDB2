package mempool

import (
    "encoding/json"
    "io/ioutil"
    "../log"
    "strconv"
    "syscall"
)

const (
    RECOVERY = true
)

var RecoveryTable = make(map[uintptr]uintptr)
var recovered = false

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

func Recovery() bool {
    if recovered {
        return false
    }
    data, err := ioutil.ReadFile(".\\malloc.json")
    if err != nil {
        return false
    }
    var mallocTable map[string]interface{}
    err = json.Unmarshal(data, &mallocTable)
    if err != nil {
        return false
    }
    newmap := make(map[uintptr]fileImage)
    for k, v := range mallocTable {
        kn, _ := strconv.Atoi(k)
        vn := v.(map[string]interface{})
        fh, ih, ip := LoadFileImage(vn["FileName"].(string))
        newmap[ip] = fileImage {
            FileHandle:     fh,
            FileName:       vn["FileName"].(string),
            ImageHandle:    ih,
        }
        RecoveryTable[uintptr(kn)] = ip
        log.WriteLog("sys", "Recovery from " + strconv.Itoa(int(kn)) + " to " + strconv.Itoa(int(ip)))
    }
    MallocTable = newmap
    return true
}

func LoadFileImage(file string) (fileHandle, imageHandle syscall.Handle, ip uintptr) {
    fileHandle = CreateFile(file, OPEN_ALWAYS)
    imageHandle = CreateFileMapping(fileHandle, 0, 0, file)
    ip = MapViewOfFile(imageHandle, 0)
    //Use the sync File IO appears to make OS refresh map view
    log.WriteLogSync("sys", "MapViewOfFile return:" + strconv.Itoa(int(ip)))
    return 
}