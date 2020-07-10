package vss

/*
#cgo CFLAGS: -I${SRCDIR}/libvshadow_prebuild/libbfio/
#cgo CFLAGS: -I${SRCDIR}/libvshadow_prebuild/common/
#cgo CFLAGS: -I${SRCDIR}/libvshadow_prebuild/libcdata/
#cgo CFLAGS: -I${SRCDIR}/libvshadow_prebuild/libcerror/
#cgo CFLAGS: -I${SRCDIR}/libvshadow_prebuild/libvshadow/
#cgo CFLAGS: -I${SRCDIR}/libvshadow_prebuild/libcfile/
#cgo CFLAGS: -I${SRCDIR}/libvshadow_prebuild/include/
#cgo CFLAGS: -Wno-discarded-qualifiers
#cgo CFLAGS: -Wno-int-to-pointer-cast
#cgo CFLAGS: -Wno-unknown-warning-option
#cgo CFLAGS: -Wno-incompatible-pointer-types-discards-qualifiers
#cgo CFLAGS: -Wno-incompatible-function-pointer-types

#cgo CFLAGS: -DHAVE_LOCAL_LIBBFIO
#cgo CFLAGS: -DHAVE_LOCAL_LIBCDATA
#cgo CFLAGS: -DHAVE_LOCAL_LIBCERROR
#cgo CFLAGS: -DHAVE_LOCAL_LIBCFILE

#cgo windows LDFLAGS: -Wl,-rpath,${SRCDIR}/libvshadow_prebuild/windows/
#cgo windows LDFLAGS: -L${SRCDIR}/libvshadow_prebuild/windows/ -lvshadow

#cgo linux LDFLAGS: -Wl,-rpath,${SRCDIR}/libvshadow_prebuild/linux/
#cgo linux LDFLAGS: -L${SRCDIR}/libvshadow_prebuild/linux/ -lvshadow

#cgo darwin LDFLAGS: -Wl,-rpath,${SRCDIR}/libvshadow_prebuild/darwin/
#cgo darwin LDFLAGS: -L${SRCDIR}/libvshadow_prebuild/macos/ -lvshadow

#include <stdlib.h>
#include <libvshadow.h>
#include <libvshadow_libcerror.h>

int libvshadow_volume_open_file_descriptor(
     libvshadow_volume_t *volume,
     int file_descriptor,
     libcerror_error_t **error);

int libvshadow_volume_open_go_reader(
     libvshadow_volume_t *volume,
     int goPointer,
     libcerror_error_t **error);
*/
import "C"
import (
	"errors"
	"io"
	"os"
	"unsafe"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/ntfs"
)

// NewByReader returns a VSS filesystem specified by a file descriptor.
func NewByReader(read io.ReaderAt) ([]*VSS, error) {
	var volume *C.libvshadow_volume_t
	var err *C.libvshadow_error_t

	if initVolume(&volume) != nil {
		return nil, errors.New("unable to initialize volume")
	}

	globalData = append(globalData, callbackData{offset: 0, reader: read})
	if C.libvshadow_volume_open_go_reader(volume, C.int(len(globalData)), &err) != 1 {
		C.libcerror_error_backtrace_fprint(err, C.stdout)
		C.libvshadow_error_free(&err)
		return nil, errors.New("unable to open volume by reader")
	}

	return getStores(volume)
}

// NewByFd returns a VSS filesystem specified by a file descriptor.
func NewByFd(fd int) ([]*VSS, error) {
	var volume *C.libvshadow_volume_t
	var err *C.libvshadow_error_t
	var fileDescriptor C.int = C.int(fd)

	if initVolume(&volume) != nil {
		return nil, errors.New("unable to initialize volume")
	}

	if C.libvshadow_volume_open_file_descriptor(volume, fileDescriptor, &err) != 1 {
		C.libcerror_error_backtrace_fprint(err, C.stdout)
		C.libvshadow_error_free(&err)
		return nil, errors.New("unable to open volume by fd")
	}

	return getStores(volume)
}

// New returns a VSS filesystem for each Volume Shadow Copy found on the system.
func NewByPath(path string) ([]*VSS, error) {
	var volume *C.libvshadow_volume_t
	var err *C.libvshadow_error_t
	var filename *C.char = C.CString(path)

	if initVolume(&volume) != nil {
		return nil, errors.New("unable to initialize volume")
	}

	if C.libvshadow_volume_open(volume, filename, C.LIBVSHADOW_OPEN_READ, &err) != 1 {
		C.libcerror_error_backtrace_fprint(err, C.stdout)
		C.libvshadow_error_free(&err)
		return nil, errors.New("unable to open volume")
	}

	return getStores(volume)
}

//getStores extract all stores found in a given volume.
func getStores(volume *C.libvshadow_volume_t) ([]*VSS, error) {
	var err *C.libvshadow_error_t
	var numberOfStores C.int = C.int(0)

	if C.libvshadow_volume_get_number_of_stores(volume, &numberOfStores, &err) != 1 {
		C.libvshadow_error_free(&err)
		return nil, errors.New("unable to get number of stores")
	}

	var stores = make([]*VSS, numberOfStores)

	for i := 0; i < int(numberOfStores); i++ {
		var store *C.libvshadow_store_t
		if C.libvshadow_volume_get_store(volume, C.int(i), &store, &err) != 1 {
			C.libvshadow_error_free(&err)
			return nil, errors.New("unable to get store")
		}
		vss := &VSS{}
		var e error
		vss.store = store
		vss.fs, e = ntfs.New(vss)
		if e != nil {
			return nil, e
		}
		stores[i] = vss
	}
	return stores, nil
}

//initVolume initializes a given volume.
func initVolume(volume **C.libvshadow_volume_t) error {
	var err *C.libvshadow_error_t

	if C.libvshadow_volume_initialize(volume, &err) != 1 {
		C.libvshadow_error_free(&err)
		return errors.New("unable to initialize volume")
	}
	return nil
}

// VSS implements a NTFS filesystem that is based on a Volume Shadow Copy.
type VSS struct {
	fs    *ntfs.FS
	store *C.libvshadow_store_t
}

// Name returns the name of the file system.
func (vss *VSS) Name() string {
	return "VSS"
}

// Open opens a file for reading.
func (vss *VSS) Open(path string) (fslib.Item, error) {
	return vss.fs.Open(path)
}

// Stat returns an os.FileInfo object that describes a file.
func (vss *VSS) Stat(path string) (os.FileInfo, error) {
	return vss.fs.Stat(path)
}

//ReadAt reads data from the internal vss store.
func (vss *VSS) ReadAt(b []byte, off int64) (n int, err error) {
	var error *C.libvshadow_error_t
	var offset Offset = Offset(off)
	var bufferLen Bufferlen = Bufferlen(len(b))

	readData := C.libvshadow_store_read_buffer_at_offset(vss.store, unsafe.Pointer(&b[0]), bufferLen, offset, &error)

	if readData == -1 {
		C.libvshadow_error_free(&error)
		return -1, errors.New("unable to read data in store")
	}

	return int(readData), nil
}
