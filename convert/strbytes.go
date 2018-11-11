package convert

import (
	"reflect"
	"unsafe"
)

// Bytes2String converts bytes to string
func Bytes2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// String2Bytes converts string to bytes
func String2Bytes(str string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
