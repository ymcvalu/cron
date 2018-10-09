package utils

import (
	"unsafe"
)

func Bytes2String(byts []byte) string {
	return *(*string)(unsafe.Pointer(&byts))
}
