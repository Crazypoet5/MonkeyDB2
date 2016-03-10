package main 
import (
    "./log"
    "time"
)

func main() {
    log.WriteLog("test", "123")
    time.Sleep(time.Second * 2)
    log.Stop()
}