package mem

import "testing"

func Test_bytes2uint(t *testing.T) {
    tmp := uintptr(7898855204580)
    r := bytes2uint(uint2bytes(tmp))
    if r != tmp {
        t.Error(r)
    }
}