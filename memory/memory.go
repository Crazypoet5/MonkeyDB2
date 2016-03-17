package memory

var ImageTable = make(map[uintptr]string)
var RecoveryTable = make(map[uintptr]uintptr)

