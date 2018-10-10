package cron

import (
	"unsafe"
)

func bytes2String(byts []byte) string {
	return *(*string)(unsafe.Pointer(&byts))
}
