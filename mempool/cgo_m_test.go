package mempool

import "testing"

func TestGetFree(t *testing.T) {
    sum := 0
    for i := 0;i < 1024;i++ {
        r := GetFree(1024 << uint(i % 14))
        r[i] = byte(i)
        sum += int(r[i])
        Release(r)
    }
    t.Error(int(sum))
    // r := Malloc(10)
    // r[5] = 5
    // k := Malloc(20)
    // k[4] = 4
    //Release(r)
}