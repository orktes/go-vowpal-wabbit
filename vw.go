package vw

// #cgo darwin pkg-config: libvw_c_wrapper
// #cgo linux pkg-config: vowpalwabbit
// #include <vowpalwabbit/vwdll.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Example a single Vowpal Wabbit example
type Example struct {
	vwHandle C.VW_HANDLE
	exHandle C.VW_EXAMPLE

	finished bool
}

// Finish finishes the example
func (ex *Example) Finish() {
	if ex.finished {
		return
	}

	C.VW_FinishExample(ex.vwHandle, ex.exHandle)
	ex.finished = true
}

// VW struct for a single Vowpal Wabbit model
type VW struct {
	handle C.VW_HANDLE

	finished bool
}

func newWithHandle(handle C.VW_HANDLE) (*VW, error) {
	vw := &VW{
		handle: handle,
	}

	return vw, nil
}

// New returns a new Vowpal Wabbit instance with given arguments
func New(args string) (*VW, error) {
	cstr := C.CString(args)
	defer C.free(unsafe.Pointer(cstr))

	return newWithHandle(C.VW_InitializeA(cstr))
}

// NewWithModelData returns new Vowpal Wabbit instance with model data
func NewWithModelData(args string, modelData []byte) (*VW, error) {
	cstr := C.CString(args)
	defer C.free(unsafe.Pointer(cstr))

	cbytes := C.CBytes(modelData)
	defer C.free(cbytes)

	return newWithHandle(C.VW_InitializeWithModel(cstr, (*C.char)(cbytes), C.size_t(len(modelData))))
}

// NewWithSeed returns a new Vowpal Wabbit instance seeded with another model
func NewWithSeed(seedModel *VW, extraArgs string) (*VW, error) {
	cstr := C.CString(extraArgs)
	defer C.free(unsafe.Pointer(cstr))

	return newWithHandle(C.VW_SeedWithModel(seedModel.handle, cstr))
}

// ReadExample parses a single text based example
func (vw *VW) ReadExample(example string) (*Example, error) {
	cstr := C.CString(example)
	defer C.free(unsafe.Pointer(cstr))

	ex := &Example{
		vwHandle: vw.handle,
		exHandle: C.VW_ReadExampleA(vw.handle, cstr),
	}

	return ex, nil
}

// Learn learns a single example and returns the score
func (vw *VW) Learn(ex *Example) float32 {
	return float32(C.VW_Learn(vw.handle, ex.exHandle))
}

// Predict returns a prediction using the example
func (vw *VW) Predict(ex *Example) float32 {
	return float32(C.VW_Predict(vw.handle, ex.exHandle))
}

// SaveModel saves model to file given during construction
func (vw *VW) SaveModel() {
	C.VW_SaveModel(vw.handle)
}

// CopyModelData returns model data as a byte slice
func (vw *VW) CopyModelData() []byte {
	var bufferHandle C.VW_IOBUF
	defer C.VW_FreeIOBuf(bufferHandle)

	var modelBytes *C.char
	var size C.size_t

	C.VW_CopyModelData(vw.handle, &bufferHandle, &modelBytes, &size)

	return C.GoBytes(unsafe.Pointer(modelBytes), C.int(size))
}

// Finish stops VW (and, eg, write weights to disk)
func (vw *VW) Finish() {
	if vw.finished {
		return
	}

	C.VW_Finish(vw.handle)
	vw.finished = true
}
