// Use of this source code is governed by a
// license that can be found in the LICENSE file.

// Package go-clamav is go wrapper for libclamav see https://docs.clamav.net/manual/Development/libclamav.html
package goclamav

/*
#include <clamav.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"
)

// Callback is used to store the interface passed to ScanFileCb. This
// object is then returned in each ClamAV callback for the duration of the
// file scan
type Callback struct {
	sync.Mutex
	nextID uintptr
	cb     map[unsafe.Pointer]interface{}
}

var callbacks = Callback{
	cb: map[unsafe.Pointer]interface{}{},
}

func setContext(i interface{}) unsafe.Pointer {
	cptr := C.malloc(1)
	if cptr == nil {
		panic("C malloc")
	}

	callbacks.Lock()
	defer callbacks.Unlock()
	callbacks.cb[cptr] = i

	return cptr
}

func findContext(key unsafe.Pointer) interface{} {
	callbacks.Lock()
	defer callbacks.Unlock()
	if v, ok := callbacks.cb[key]; ok {
		return v
	}
	return nil
}

func deleteContext(key unsafe.Pointer) error {
	callbacks.Lock()
	defer callbacks.Unlock()
	if _, ok := callbacks.cb[key]; ok {
		delete(callbacks.cb, key)
		C.free(key)
		return nil
	}
	return errors.New("no context to delete")
}

type Clamav struct {
	engine  *C.struct_cl_engine
	signo   uint
	options *C.struct_cl_scan_options
}

// Init new clamav instance
func (c *Clamav) Init(options SCAN_OPTIONS) error {
	c.engine = (*C.struct_cl_engine)(C.cl_engine_new())

	scanOptions := &C.struct_cl_scan_options{
		general:   C.uint(options.General),
		heuristic: C.uint(options.Heuristic),
		parse:     C.uint(options.Parse),
		mail:      C.uint(options.Mail),
		dev:       C.uint(options.Dev),
	}
	c.options = scanOptions

	ret := ErrorCode(C.cl_init(CL_INIT_DEFAULT))
	if ret != CL_SUCCESS {
		err := Strerr(ret)
		return err
	}
	return nil
}

// Use the CvdVerify to verify a database directly:
// As the comment block explains, this will load-test the database. Be advised
// that for some larger databases, this may use a fair bit system RAM.
func (c *Clamav) CvdVerify(path string) error {
	_, err := os.Stat(path)
	existed := !os.IsNotExist(err)
	if !existed {
		err := errors.New(fmt.Sprintf("db %s is not exists!", path))
		return err
	}

	fp := C.CString(path)
	defer C.free(unsafe.Pointer(fp))

	ret := ErrorCode(C.cl_cvdverify(fp))
	if ret != CL_SUCCESS {
		err := Strerr(ret)
		return err
	}

	return nil
}

// Load clamav virus database
func (c *Clamav) LoadDB(path string, dbopts uint) (uint, error) {
	_, err := os.Stat(path)
	existed := !os.IsNotExist(err)
	if !existed {
		err := errors.New(fmt.Sprintf("db %s is not exists!", path))
		return 0, err
	}

	var signo uint
	fp := C.CString(path)
	defer C.free(unsafe.Pointer(fp))

	ret := ErrorCode(C.cl_load(fp, (*C.struct_cl_engine)(c.engine), (*C.uint)(unsafe.Pointer(&signo)), C.uint(dbopts)))
	if ret != CL_SUCCESS {
		err := Strerr(ret)
		return 0, err
	}

	return signo, nil
}

// When all required databases are loaded you should prepare the detection engine by calling CompileEngine
func (c *Clamav) CompileEngine() error {
	ret := ErrorCode(C.cl_engine_compile((*C.struct_cl_engine)(c.engine)))
	if ret != CL_SUCCESS {
		err := Strerr(ret)
		return err
	}

	return nil
}

// EngineSetNum sets a number in the specified field of the engine configuration.
// Certain fields accept only 32-bit numbers, silently truncating the higher bits
// of the engine config. See dat.go for more information.
func (c *Clamav) EngineSetNum(field EngineField, num uint64) error {
	ret := ErrorCode(C.cl_engine_set_num((*C.struct_cl_engine)(c.engine), C.enum_cl_engine_field(field), C.longlong(num)))
	if ret != CL_SUCCESS {
		err := Strerr(ret)
		return err
	}

	return nil
}

// EngineGetNum acquires a number from the specified field of the engine configuration. Tests show that
// the ClamAV library will not overflow 32-bit fields, so a GetNum on a 32-bit field can safely be
// cast to uint32.
func (c *Clamav) EngineGetNum(field EngineField) (uint64, error) {
	var ret ErrorCode
	ne := (*C.struct_cl_engine)(c.engine)
	num := uint64(C.cl_engine_get_num(ne, C.enum_cl_engine_field(field), (*C.int)(unsafe.Pointer(&ret))))
	if ret != CL_SUCCESS {
		err := Strerr(ret)
		return num, err
	}

	return num, nil
}

// Free the memory allocated to clamav instance, Free should be called
// when the engine is no longer in use.
func (c *Clamav) Free() error {
	ret := ErrorCode(C.cl_engine_free((*C.struct_cl_engine)(c.engine)))
	if ret == CL_SUCCESS {
		return nil
	}
	return Strerr(ret)
}

// ScanMapCB scans custom data
func (c *Clamav) ScanMapCB(fmap *Fmap, fileName string, context interface{}) (uint, string, error) {
	var scanned C.ulong
	var virusName *C.char

	fn := C.CString(fileName)
	defer C.free(unsafe.Pointer(fn))

	// find where to store the context in our callback map. we do _not_ pass the context to
	// C directly because aggressive garbage collection will move it around
	ctx := setContext(context)
	// cleanup
	defer deleteContext(ctx)

	ret := ErrorCode(C.cl_scanmap_callback((*C.cl_fmap_t)(fmap), fn, &virusName, &scanned, (*C.struct_cl_engine)(c.engine), (*C.struct_cl_scan_options)(c.options), unsafe.Pointer(ctx)))
	defer CloseMemory(fmap)
	// clean
	if ret == CL_SUCCESS {
		return uint(scanned), "", nil
	}
	// virus
	if ret == CL_VIRUS {
		return uint(scanned), C.GoString(virusName), Strerr(ret)
	}
	// error
	return 0, "", Strerr(ret)
}

// ScanFile scans a single file for viruses using the ClamAV databases. It returns the number of bytes
// read from the file (if found), the virus name and an error code.
// If the file is clean, the virus name is empty and the error code is nil,but if the file is insecure, the
// error code is "Virus(es) detected" and virus name is the matching rule.
func (c *Clamav) ScanFile(path string) (uint, string, error) {
	fp := C.CString(path)
	defer C.free(unsafe.Pointer(fp))

	var virusName *C.char
	var scanned C.ulong

	ret := ErrorCode(C.cl_scanfile(fp, &virusName, &scanned, (*C.struct_cl_engine)(c.engine), (*C.struct_cl_scan_options)(c.options)))
	// clean
	if ret == CL_SUCCESS {
		return uint(scanned), "", nil
	}
	// virus
	if ret == CL_VIRUS {
		return uint(scanned), C.GoString(virusName), Strerr(ret)
	}
	// error
	return 0, "", Strerr(ret)
}

// ScanFileCB scans a single file for viruses using the ClamAV databases and using callbacks from
// ClamAV to read/resolve file data. The callbacks can be used to scan files in memory, to scan multiple
// files inside archives, etc. The function returns the number of bytes
// read from the file (if found), the virus name and an error code.
// If the file is clean, the virus name is empty and the error code is nil,but if the file is insecure, the
// error code is "Virus(es) detected" and virus name is the matching rule.
// The context argument will be sent back to the callbacks, so effort must be made to retain it
// throughout the execution of the scan from garbage collection
func (c *Clamav) ScanFileCB(path string, context interface{}) (uint, string, error) {
	fp := C.CString(path)
	defer C.free(unsafe.Pointer(fp))

	// find where to store the context in our callback map. we do _not_ pass the context to
	// C directly because aggressive garbage collection will move it around
	ctx := setContext(context)
	// cleanup
	defer deleteContext(ctx)

	var virusName *C.char
	var scanned C.ulong

	ret := ErrorCode(C.cl_scanfile_callback(fp, &virusName, &scanned, (*C.struct_cl_engine)(c.engine), (*C.struct_cl_scan_options)(c.options), ctx))
	// clean
	if ret == CL_SUCCESS {
		return 0, "", nil
	}
	// virus
	if ret == CL_VIRUS {
		return uint(scanned), C.GoString(virusName), Strerr(ret)
	}
	// error
	return 0, "", Strerr(ret)
}

// ScanDesc scans a file descriptor for viruses using the ClamAV databases. It returns the number of bytes
// read from the file (if found), the virus name and an error code.
// If the file is clean, the virus name is empty and the error code is nil,but if the file is insecure, the
// error code is "Virus(es) detected" and virus name is the matching rule.
func (c *Clamav) ScanDesc(desc int32, fileName string) (uint, string, error) {
	var scanned C.ulong
	var virusName *C.char

	fn := C.CString(fileName)
	defer C.free(unsafe.Pointer(fn))

	ret := ErrorCode(C.cl_scandesc(C.int(desc), fn, &virusName, &scanned, (*C.struct_cl_engine)(c.engine), (*C.struct_cl_scan_options)(c.options)))
	// clean
	if ret == CL_SUCCESS {
		return 0, "", nil
	}
	// virus
	if ret == CL_VIRUS {
		return uint(scanned), C.GoString(virusName), Strerr(ret)
	}
	// error
	return 0, "", Strerr(ret)
}

// OpenMemory creates an object from the given memory that can be scanned using ScanMapCb
func OpenMemory(start []byte) *Fmap {
	return (*Fmap)(C.cl_fmap_open_memory(unsafe.Pointer(&start[0]), C.size_t(len(start))))
}

// CloseMemory destroys the fmap associated with an in-memory object
func CloseMemory(f *Fmap) {
	C.cl_fmap_close((*C.cl_fmap_t)(f))
}
