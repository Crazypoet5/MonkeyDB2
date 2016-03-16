package mempool

import (
    "unsafe"
    "reflect"
    "./heap"
    "time"
    "syscall"
)

const MAX_LEVEL = 15        //It means we can use 1k, 2k, ... , 2^14 * 1k

const MAX_BLOCKS_LIMIT = 2  //When free max blocks more than MAX_BLOCKS_LIMIT, when release memory, we will free them  

type blocks [][]byte

func (b *blocks) Swap(i, j int) {
    header_i := (*reflect.SliceHeader)(unsafe.Pointer(&(*b)[i]))
    header_j := (*reflect.SliceHeader)(unsafe.Pointer(&(*b)[j]))
    var header_t reflect.SliceHeader
    assign(&header_t, *header_i)
    assign(header_i, *header_j)
    assign(header_j, header_t)
}

func (b *blocks) Len() int {
    return len(*b)
}

func (b *blocks) Less(i, j int) bool {
    header_i := (*reflect.SliceHeader)(unsafe.Pointer(&(*b)[i]))
    header_j := (*reflect.SliceHeader)(unsafe.Pointer(&(*b)[j]))
    return header_i.Data < header_j.Data
}

func (b *blocks) Push(x interface{}) {
    if bs, ok := x.([]byte);ok {
        *b = append(*b, bs)
    }
}

func (b *blocks) Pop() interface{} {
    p := (*b)[len(*b) - 1]
    *b = (*b)[:len(*b) - 1]
    return p
}

func assign(header_t *reflect.SliceHeader, header_i reflect.SliceHeader) {
    header_t.Data = header_i.Data
    header_t.Len = header_i.Len
    header_t.Cap = header_i.Cap
}

type Pool struct {
    list        [MAX_LEVEL]blocks
}

var pool Pool

type fileImage struct {
    FileHandle  syscall.Handle
    ImageHandle syscall.Handle
    FileName    string
}


func init() {
    if RECOVERY && Recovery() {
        return
    }
    for i := 0;i < MAX_LEVEL;i++ {
        heap.Push(&pool.list[i], Malloc(1024 << uint(i)))
        time.Sleep(time.Second)
    }
}

func canMerge(b1, b2 []byte) bool {
    header1 := (*reflect.SliceHeader)(unsafe.Pointer(&b1))
    header2 := (*reflect.SliceHeader)(unsafe.Pointer(&b2))
    return int(header1.Data) + header1.Cap == int(header2.Data) ||
     int(header2.Data) + header2.Cap == int(header2.Data)
}

func merge(prev, back []byte) []byte {
    var header reflect.SliceHeader
    header1 := (*reflect.SliceHeader)(unsafe.Pointer(&prev))
    header2 := (*reflect.SliceHeader)(unsafe.Pointer(&back))
    header.Cap = header1.Cap + header2.Cap
    header.Len = header.Cap
    header.Data = header1.Data
    return *(*[]byte)(unsafe.Pointer(&header))
}

/*
    InsertFree function insert a free space buffer to list[n]
*/
func InsertFree(n int, b []byte) {
    if n >= MAX_LEVEL {
        return
    }
    // Max blocks remain little num
    if n == MAX_LEVEL - 1 && pool.list[MAX_LEVEL - 1].Len() >= MAX_BLOCKS_LIMIT {
        Free(b)
        return
    }
    //Merge
    p := heap.Push(&pool.list[n], b)
    if p > 0 && canMerge(pool.list[n][p - 1], pool.list[n][p]) {
        big := merge(pool.list[n][p - 1], pool.list[n][p])
        heap.Remove(&pool.list[n], p)
        InsertFree(n + 1, big)
    } else if p < pool.list[n].Len() - 1 && canMerge(pool.list[n][p], pool.list[n][p + 1]) {
        big := merge(pool.list[n][p], pool.list[n][p + 1])
        heap.Remove(&pool.list[n], p)
        InsertFree(n + 1, big)
    }
    
}

func slice(full []byte) ([]byte, []byte) {
    header := (*reflect.SliceHeader)(unsafe.Pointer(&full))
    var header1, header2 reflect.SliceHeader
    header1.Cap = header.Cap / 2
    header1.Len = header1.Cap
    header1.Data = header.Data
    header2.Cap = header1.Cap
    header2.Len = header2.Cap
    header2.Data = uintptr(int(header1.Data) + header1.Len)
    return *(*[]byte)(unsafe.Pointer(&header1)),
    *(*[]byte)(unsafe.Pointer(&header2))
}

/*
    Slice function slice a big block of n to 2 parts, then insert one to list[n-1] and return another
*/
func Slice(n int) []byte {
    if n > MAX_LEVEL || n <= 0 {
        panic("Unexpected block applied to slice!")
    }
    if n == MAX_LEVEL {
        return Malloc(1024 << (MAX_LEVEL - 1))
    }
    full := getFree(n)
    half1, half2 := slice(full)
    InsertFree(n - 1, half2)
    return half1
}

func getFree(n int) []byte {
    if pool.list[n].Len() != 0 {
        return heap.Pop(&pool.list[n]).([]byte)
    }
    return Slice(n + 1)
}

func GetFree(size int) []byte {
    n := size >> 10
    i := 0
    for n > 0 {
        n >>= 1
        i++
    }
    return getFree(i)
}

// Warning: DONOT release a []byte created by Go
func Release(buffer []byte) {
    header := (*reflect.SliceHeader)(unsafe.Pointer(&buffer))
    n := header.Cap >> 10
    i := 0
    for n > 0 {
        n >>= 1
        i++
    }
    InsertFree(i, buffer)
}