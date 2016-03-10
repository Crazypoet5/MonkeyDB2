package log

import (
    "os"
    "fmt"
    "time"
    "../common"
)

type msg struct {
    kind string
    content interface{}
    abort   bool
}

var queue = make(chan msg, 10)

func checkFileIsExist(filename string) (bool) {
    var exist = true;
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exist = false;
    }
    return exist;
}

func writeLog(kind string, content interface{}) {
    var f *os.File
    var err error
    filename := common.FixPath(common.GetCurrentDirectory() + "\\log\\" + kind + ".log")
    fmt.Println(filename)
    if checkFileIsExist(filename) {
        f, err = os.OpenFile(filename, os.O_APPEND, 0666)
    }else {
        f, err = os.Create(filename)
    }
    if err != nil {
        panic("open " + filename + " fail!")
    }
    defer f.Close()
    datetime := time.Now().String()
    fmt.Fprintf(f, "[ %s ]", datetime)
    fmt.Fprint(f, content)
}

func run() {
    for {
        m := <- queue
        fmt.Println(m.kind)
        if m.abort {
            return
        }
        writeLog(m.kind, m.content)
    }
}

func init() {
    go run()
}

func WriteLog(kind string, content interface{}) {
    queue <- msg {
        kind:   kind,
        content:    content,
        abort:  false,
    }
}

func Stop() {
    queue <- msg {
        abort: true,
    }
}
//TODO: Test