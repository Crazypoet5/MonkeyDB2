package memory

//#include <stdlib.h>
import "C"

import (
	"container/list"
	"os"
	"strconv"
	"unsafe"

	"../common"
)

var DataBlockTable = make(map[uintptr]*DataBlock)
var DataBlockList list.List
var commonBuffer = make([]byte, 1024*1024) // To clear the file quickly
var count = 0                              // Windows only support 100ns level

// CreateImage creates a image file and returns the address
func CreateImage(size int) (ip *DataBlock, err error) {
	defer SignalBackup()
	filename := common.COMMON_DIR + "\\image\\" + strconv.Itoa(count)
	count++
	ip = &DataBlock{
		RawPtr: uintptr(C.malloc(C.size_t(size))),
		Size:   uint(size),
	}
	file, err := os.Create(filename)
	defer file.Close()
	for i := size; i > 0; i -= 1024 * 1024 {
		if i < 1024*1024 {
			file.Write(commonBuffer[:i])
		} else {
			file.Write(commonBuffer)
		}
	}
	ImageTable[ip.RawPtr] = filename
	DataBlockTable[ip.RawPtr] = ip
	DataBlockList.PushBack(ip)
	return
}

// ReallocImage creates a new bigger image file and returns the new address with copying data
func ReallocImage(ip *DataBlock, size int) (*DataBlock, error) {
	defer SignalBackup()
	filename := common.COMMON_DIR + "\\image\\" + strconv.Itoa(count)
	count++
	os.Remove(ImageTable[ip.RawPtr])
	ipNew := &DataBlock{
		RawPtr: uintptr(C.malloc(C.size_t(size))),
		Size:   uint(size),
	}
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	for i := size; i > 0; i -= 1024 * 1024 {
		if i < 1024*1024 {
			file.Write(commonBuffer[:i])
		} else {
			file.Write(commonBuffer)
		}
	}
	Copy(ipNew, ip)
	delete(ImageTable, ip.RawPtr)
	C.free(unsafe.Pointer(ip.RawPtr))
	ImageTable[ipNew.RawPtr] = filename
	DataBlockTable[ipNew.RawPtr] = ip
	RemoveBlock(ip)
	DataBlockList.PushBack(ipNew)
	return ipNew, nil
}

// ReleaseImage release the image and delete its related file
func ReleaseImage(ip *DataBlock) (err error) {
	defer SignalBackup()
	_, ok := ImageTable[ip.RawPtr]
	if !ok {
		err = NOT_FOUND_ADDRESS
		return
	}
	err = os.Remove(ImageTable[ip.RawPtr])
	RemoveBlock(ip)
	C.free(unsafe.Pointer(ip.RawPtr))
	delete(ImageTable, ip.RawPtr)
	delete(DataBlockTable, ip.RawPtr)
	for l := DataBlockList.Front(); l != nil; l = l.Next() {
		if l.Value.(*DataBlock) == ip {
			DataBlockList.Remove(l)
		}
	}
	return
}

func RemoveBlock(ip *DataBlock) {
	var l *list.Element
	for l = DataBlockList.Front(); l != nil; l = l.Next() {
		if l.Value.(*DataBlock) == ip {
			break
		}
	}
	DataBlockList.Remove(l)
	delete(DataBlockTable, ip.RawPtr)
}
