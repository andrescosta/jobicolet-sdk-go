package sdk

// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

type EventFunc func(string) (uint64, string)

var OnEvent EventFunc

type Level uint32

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
)

const (
	NoError uint64 = 0
)

func Log(level Level, message string) {
	ptr, size := StringToPtr(message)
	_log(level, ptr, size)
	runtime.KeepAlive(message) // keep message alive until ptr is no longer needed.
}

//go:wasmimport env log
func _log(level Level, ptr, size uint32)

func PtrToString(ptr uint32, size uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

func StringToPtr(s string) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

func StringToLeakedPtr(s string) (uint32, uint32) {
	size := C.ulong(len(s))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), s)
	return uint32(uintptr(ptr)), uint32(size)
}

//export ver
func ver() uint32 {
	return 0
}

//export event
func event(ptr, size uint32) (uint64, uint64) {
	data := PtrToString(ptr, size)
	errno, result := OnEvent(data)
	ptr, size = StringToLeakedPtr(result)
	ptrRes := (uint64(ptr) << uint64(32)) | uint64(size)
	return errno, ptrRes
}
