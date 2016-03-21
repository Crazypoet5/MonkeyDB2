package mempool

import "testing"
import "time"

func TestGetFree(t *testing.T) {
    sum := 0
    for i := 0;i < 100;i++ {
        r := GetFree(64 << uint(i % 14))
        r[25] = byte(i)
        sum += int(r[25])
        time.Sleep(time.Millisecond)    //If too fast then it will failed!
        Release(r)
    }
    t.Error(int(sum))
    // r := Malloc(10)
    // r[5] = 5
    // k := Malloc(20)
    // k[4] = 4
    //Release(r)
}