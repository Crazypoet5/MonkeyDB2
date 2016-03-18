package memory

import (
    "io/ioutil"
)

var startBackup = make(chan bool)

func (b *DataBlock) SyncToFile() error {
    data, err := b.Read(0, b.Size)
    if err != nil {
        return err
    }
    filename, ok := ImageTable[b.RawPtr]
    if !ok {
        return NOT_FOUND_ADDRESS
    }
    return ioutil.WriteFile(filename, data, 0666)
}

func SyncAllImageToFile() {
    for l := DataBlockList.Front();l != nil;l = l.Next() {
        b, ok := l.Value.(*DataBlock)
        if !ok {
            continue
        }
        data, _ := b.Read(0, b.Size)
        filename, ok := ImageTable[b.RawPtr]
        if !ok {
            continue
        }
        ioutil.WriteFile(filename, data, 0666)
    }
}

func BackupRoutine() {
    
}