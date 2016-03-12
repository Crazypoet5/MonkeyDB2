package mempool

import "testing"

func TestGetFree(t *testing.T) {
    // for i := 0;i < 1024;i++ {
    //     r := GetFree(1024 << uint(i % 15))
    //     Release(r)
    // }
    r := Malloc(10)
    r[5] = 5
    k := Malloc(20)
    k[4] = 4
    //Release(r)
}