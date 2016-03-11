package csbt

import "testing"

func expect(exp int, real int, t *testing.T) {
    if real != exp {
        t.Error("Expected:", exp, "Get:", real)
    }
}

func Test_binarySearch(t *testing.T) {
    array := []uint{1, 3, 5, 7, 9, 11, 13, 15, 17}
    i, _ := binarySearch(0, array, 0, 8)
    expect(-1, i, t)
    i, _ = binarySearch(18, array, 0, 8)
    expect(9, i, t)
    i, _ = binarySearch(10, array, 0, 8)
    expect(4, i, t)
}