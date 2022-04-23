package vss

/*
#include <libvshadow_libcerror.h>
#include <stdio.h>

void *memcpy(void *dest, const void * src, size_t n);
*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

type callbackData struct {
	reader io.ReaderAt
	offset int
}

var globalData []callbackData

//export reader_handle_free
func reader_handle_free(handle **C.intptr_t, err **C.libcerror_error_t) C.int {
	fmt.Println("Handle free not implemented...")
	return C.int(1)
}

//export reader_handle_clone
func reader_handle_clone(src **C.intptr_t, dst *C.intptr_t, err **C.libcerror_error_t) C.int {
	fmt.Println("Handle clone not implemented...")
	return C.int(-1)
}

//export reader_handle_open
func reader_handle_open(handle *C.intptr_t, flags C.int, err **C.libcerror_error_t) C.int {
	fmt.Println("Handle open not implemented...")
	return C.int(1)
}

//export reader_handle_close
func reader_handle_close(handle *C.intptr_t, err **C.libcerror_error_t) C.int {
	fmt.Println("Handle close not implemented...")
	return C.int(1)
}

//export reader_handle_read
func reader_handle_read(handle *C.intptr_t, buffer *C.uint8_t, size C.size_t, err **C.libcerror_error_t) C.ssize_t {
	/*var test *VSSReader = (*VSSReader)(unsafe.Pointer(handle))
	var reader *io.ReaderAt = (*io.ReaderAt)(unsafe.Pointer(test.reader))*/
	index := int(uintptr(unsafe.Pointer(handle)))
	reader := globalData[index-1].reader
	offset := globalData[index-1].offset

	buf := make([]byte, int(size))
	reader.ReadAt(buf, int64(offset)) // nolint: errcheck

	C.memcpy(unsafe.Pointer(buffer), unsafe.Pointer(&(buf[0])), size)

	return C.ssize_t(size)
}

//export reader_handle_write
func reader_handle_write(handle *C.intptr_t, buffer *C.uint8_t, size C.size_t, err **C.libcerror_error_t) C.ssize_t {
	fmt.Println("Handle write not implemented...")
	return C.ssize_t(0)
}

//export reader_handle_seak_offset
func reader_handle_seak_offset(handle *C.intptr_t, off C.off64_t, whence C.int, err **C.libcerror_error_t) C.ssize_t {
	// var test *VSSReader = (*VSSReader)(unsafe.Pointer(handle))
	// test.offset = int(off)
	index := int(uintptr(unsafe.Pointer(handle)))
	globalData[index-1].offset = int(off)
	return C.ssize_t(off) // nolint: unconvert
}

//export reader_handle_exists
func reader_handle_exists(handle *C.intptr_t, err **C.libcerror_error_t) C.int {
	fmt.Println("Handle exists not implemented...")
	return C.int(-1)
}

//export reader_handle_is_open
func reader_handle_is_open(handle *C.intptr_t, err **C.libcerror_error_t) C.int {
	return C.int(1)
}

//export reader_handle_get_size
func reader_handle_get_size(handle *C.intptr_t, size *C.size64_t, err **C.libcerror_error_t) C.int {
	fmt.Println("Handle get_size not implemented...")
	return C.int(-1)
}
