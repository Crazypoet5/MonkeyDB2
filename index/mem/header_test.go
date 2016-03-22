package mem

import "testing"

func Test_bytes2int(t *testing.T) {
    tmp := uint(7898855204580)
    r := bytes2uint(uint2bytes(tmp))
    if r != tmp {
        t.Error(r)
    }
}