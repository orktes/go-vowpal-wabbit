package vw

// #cgo CXXFLAGS: -std=c++11 -I${SRCDIR}/extra -I${SRCDIR} -O3 -Wall -g -Wno-sign-compare -Wno-unused-function -I/Library/Developer/CommandLineTools/usr/include/c++/v1 -I/usr/local/include
// #cgo LDFLAGS: -lstdc++
// #cgo pkg-config: libvw_c_wrapper libvw
// #include <vowpalwabbit/vwdll.h>
// #include <stdlib.h>
// #include "lib.hpp"
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

// GetLabel returns label
func (ex *Example) GetLabel() float32 {
	return float32(C.VW_GetLabel(ex.exHandle))
}

// GetInitial returns initial
func (ex *Example) GetInitial() float32 {
	return float32(C.VW_GetInitial(ex.exHandle))
}

// GetImportance returns importance
func (ex *Example) GetImportance() float32 {
	return float32(C.VW_GetImportance(ex.exHandle))
}

// GetPrediction returns prediction
func (ex *Example) GetPrediction() float32 {
	return float32(C.VW_GetPrediction(ex.exHandle))
}

// GetCostSensitivePrediction returns initial
func (ex *Example) GetCostSensitivePrediction() float32 {
	return float32(C.VW_GetCostSensitivePrediction(ex.exHandle))
}

// GetTopicPrediction returns topic prediction for index
func (ex *Example) GetTopicPrediction(i int) float32 {
	return float32(C.VW_GetTopicPrediction(ex.exHandle, C.size_t(i)))
}

// GetActionScore returns action score for index
func (ex *Example) GetActionScore(i int) float32 {
	return float32(C.VW_GetActionScore(ex.exHandle, C.size_t(i)))
}

// GetActionScores returns all action scores as a slice
func (ex *Example) GetActionScores() []float32 {
	length := ex.GetActionScoreLength()
	if length < 0 {
		return nil
	}

	scores := make([]float32, length)

	for i := 0; i < length; i++ {
		scores[i] = ex.GetActionScore(i)
	}

	return scores
}

// GetActionScoreLength get action score length
func (ex *Example) GetActionScoreLength() int {
	return int(C.VW_GetActionScoreLength(ex.exHandle))
}

// GetTagLength returns tag length
func (ex *Example) GetTagLength() int {
	return int(C.VW_GetTagLength(ex.exHandle))
}

// GetTag returns tag
func (ex *Example) GetTag() string {
	return C.GoString(C.VW_GetTag(ex.exHandle))
}

// GetFeatureNumber returns number of features
func (ex *Example) GetFeatureNumber() int {
	return int(C.VW_GetFeatureNumber(ex.exHandle))
}

// GetConfidence returns confidence
func (ex *Example) GetConfidence(i int) float32 {
	return float32(C.VW_GetConfidence(ex.exHandle))
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

// ReadDecisionServiceJSON reads examples from Decision Service JSON format
func (vw *VW) ReadDecisionServiceJSON(json string) ([]*Example, error) {
	return nil, nil
}

// Learn learns a single example and returns the score
func (vw *VW) Learn(ex *Example) float32 {
	return float32(C.VW_Learn(vw.handle, ex.exHandle))
}

// Predict returns a prediction using the example
func (vw *VW) Predict(ex *Example) float32 {
	return float32(C.VW_Predict(vw.handle, ex.exHandle))
}

// PredictCostSensitive returns a cost sensitive prediction using the example
func (vw *VW) PredictCostSensitive(ex *Example) float32 {
	return float32(C.VW_PredictCostSensitive(vw.handle, ex.exHandle))
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
