// Package vw implements golang bindings for Vowpal Wabbit Online Learning library.
// API is designed to closely resemble the C API of Vowpal Wabbit with minor changes
// done to make it more convinient to be used from go. Library is not thread safe
// and additional locking is required if libary is being called from multiple goroutines.
package vw

// #cgo CXXFLAGS: -std=c++11 -I${SRCDIR}/extra -I${SRCDIR} -O3 -Wall -g -Wno-sign-compare -Wno-unused-function -I/Library/Developer/CommandLineTools/usr/include/c++/v1 -I/usr/local/include
// #cgo LDFLAGS: -lstdc++
// #cgo pkg-config: libvw_c_wrapper libvw RapidJSON
// #include <vowpalwabbit/vwdll.h>
// #include <stdlib.h>
// #include "lib.hpp"
import "C"
import (
	"runtime"
	"unsafe"
)

type vwError = C.VW_ERROR

// VW struct for a single Vowpal Wabbit model
type VW struct {
	handle      C.VW_HANDLE
	examplePool C.VW_EXAMPLE_POOL_HANDLE

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

func (vw *VW) getExamplePool() C.VW_EXAMPLE_POOL_HANDLE {
	if vw.examplePool == nil {
		vw.examplePool = C.VW_CreateExamplePool(vw.handle)
	}

	return vw.examplePool
}

// ReadExample parses a single text based example
func (vw *VW) ReadExample(example string) (*Example, error) {
	cstr := C.CString(example)
	defer C.free(unsafe.Pointer(cstr))

	ex := &Example{
		vw:       vw,
		exHandle: C.VW_ReadExampleA(vw.handle, cstr),
	}

	return ex, nil
}

// ReadDecisionServiceJSON reads examples from Decision Service JSON format
func (vw *VW) ReadDecisionServiceJSON(json string) (ExampleList, error) {
	bytes := make([]byte, len(json)+1)
	copy(bytes, json)

	// Terminate bytes
	bytes[len(bytes)-1] = 0

	return vw.ReadDecisionServiceJSONFromBytes(bytes)
}

// ReadDecisionServiceJSONFromBytes like ReadDecisionServiceJSON but takes in a byte slice (need to be null terminated)
func (vw *VW) ReadDecisionServiceJSONFromBytes(json []byte) (ExampleList, error) {
	jsonPtr := (*C.char)(unsafe.Pointer(&json[0]))

	var size C.ulong
	var vwerr vwError

	pool := vw.getExamplePool()
	examplePtr := C.VW_ReadDSJSONExampleSafe(vw.handle, pool, jsonPtr, &size, &vwerr)
	if err := checkError(vwerr); err != nil {
		return nil, err
	}

	exampleHandles := (*[1 << 30]C.VW_EXAMPLE)(unsafe.Pointer(examplePtr))[:size:size]
	examples := make([]*Example, size)

	for i, handle := range exampleHandles {
		examples[i] = &Example{
			vw:       vw,
			exHandle: handle,
			fromPool: true,
		}
	}

	C.free(unsafe.Pointer(examplePtr))
	runtime.KeepAlive(json)

	return ExampleList(examples), nil
}

// ReadJSON reads examples from format
func (vw *VW) ReadJSON(json string) (ExampleList, error) {
	cstr := C.CString(json)
	defer C.free(unsafe.Pointer(cstr))

	var size C.ulong
	var vwerr vwError

	pool := vw.getExamplePool()
	examplePtr := C.VW_ReadJSONExampleSafe(vw.handle, pool, cstr, &size, &vwerr)
	if err := checkError(vwerr); err != nil {
		return nil, err
	}

	exampleHandles := (*[1 << 30]C.VW_EXAMPLE)(unsafe.Pointer(examplePtr))[:size:size]
	examples := make([]*Example, size)

	for i, handle := range exampleHandles {
		examples[i] = &Example{
			exHandle: handle,
			vw:       vw,
			fromPool: true,
		}
	}

	C.free(unsafe.Pointer(examplePtr))

	return ExampleList(examples), nil
}

// Learn learns a single example and returns the score
func (vw *VW) Learn(ex *Example) float32 {
	return float32(C.VW_Learn(vw.handle, ex.exHandle))
}

// Predict returns a prediction using the example
func (vw *VW) Predict(ex *Example) float32 {
	return float32(C.VW_Predict(vw.handle, ex.exHandle))
}

// MultiLineLearn learns from a list of example and returns the score
func (vw *VW) MultiLineLearn(exs []*Example) error {
	if len(exs) == 0 {
		return nil
	}

	exHandles := make([]C.VW_EXAMPLE, len(exs))
	for i, ex := range exs {
		exHandles[i] = ex.exHandle
	}

	exPtr := (*C.VW_EXAMPLE)(unsafe.Pointer(&exHandles[0]))

	var vwerr vwError
	C.VW_MultiLineLearnSafe(vw.handle, exPtr, C.size_t(len(exHandles)), &vwerr)
	if err := checkError(vwerr); err != nil {
		return err
	}

	runtime.KeepAlive(exHandles)
	return nil
}

// MultiLinePredict predic using a list of examples
func (vw *VW) MultiLinePredict(exs []*Example) error {
	if len(exs) == 0 {
		return nil
	}

	exHandles := make([]C.VW_EXAMPLE, len(exs))
	for i, ex := range exs {
		exHandles[i] = ex.exHandle
	}

	exPtr := (*C.VW_EXAMPLE)(unsafe.Pointer(&exHandles[0]))

	var vwerr vwError
	C.VW_MultiLinePredictSafe(vw.handle, exPtr, C.size_t(len(exHandles)), &vwerr)
	if err := checkError(vwerr); err != nil {
		return err
	}

	runtime.KeepAlive(exHandles)
	return nil
}

// GetLearningRate returns the current learning rate
func (vw *VW) GetLearningRate() (float32, error) {
	var vwerr vwError
	learningRate := C.VW_GetLearningRate(vw.handle, &vwerr)
	if err := checkError(vwerr); err != nil {
		return 0.0, err
	}

	return float32(learningRate), nil
}

// PredictCostSensitive returns a cost sensitive prediction using the example
func (vw *VW) PredictCostSensitive(ex *Example) float32 {
	return float32(C.VW_PredictCostSensitive(vw.handle, ex.exHandle))
}

// StartParser starts the parser
func (vw *VW) StartParser() error {
	C.VW_StartParser(vw.handle)
	return nil
}

// EndParser ends the parser
func (vw *VW) EndParser() error {
	C.VW_EndParser(vw.handle)
	return nil
}

// SaveModel saves model to file given during construction
func (vw *VW) SaveModel() {
	C.VW_SaveModel(vw.handle)
}

// PerformanceStatistics writes current performance stats
func (vw *VW) PerformanceStatistics() (PerformanceStatistics, error) {
	var stats PerformanceStatistics

	var vwerr vwError
	statsCStruct := C.VW_PerformanceStats(vw.handle, &vwerr)
	if err := checkError(vwerr); err != nil {
		return stats, err
	}

	stats.CurrentPass = uint64(statsCStruct.current_pass)
	stats.NumberOfFeatures = uint64(statsCStruct.number_of_features)
	stats.NumberOfExamples = uint64(statsCStruct.number_of_examples)
	stats.WeightedExampleSum = float64(statsCStruct.weighted_example_sum)
	stats.WeightedLabelSum = float64(statsCStruct.weighted_label_sum)
	stats.AverageLoss = float64(statsCStruct.average_loss)
	stats.BestConstant = float64(statsCStruct.best_constant)
	stats.BestConstantLoss = float64(statsCStruct.best_constant_loss)

	return stats, nil
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

	if vw.examplePool != nil {
		C.VW_ReleaseExamplePool(vw.examplePool)
		vw.examplePool = nil
	}

	C.VW_Finish(vw.handle)
	vw.handle = nil
	vw.finished = true
}

// EndOfPass ends current pass and syncs metrics
func (vw *VW) EndOfPass() error {
	var vwerr vwError
	C.VW_EndOfPass(vw.handle, &vwerr)
	if err := checkError(vwerr); err != nil {
		return err
	}
	return nil
}

// FinishPasses finalize passes
func (vw *VW) FinishPasses() error {
	C.VW_Finish_Passes(vw.handle)
	return nil
}

// SyncStats syncs stats from the learner (should be called before vw.PerformanceStatistics)
func (vw *VW) SyncStats() error {
	var vwerr vwError
	C.VW_SyncStats(vw.handle, &vwerr)
	if err := checkError(vwerr); err != nil {
		return err
	}
	return nil
}

// Example a single Vowpal Wabbit example
type Example struct {
	exHandle C.VW_EXAMPLE
	vw       *VW
	fromPool bool

	finished bool
}

// Finish finishes the example
func (ex *Example) Finish() {
	if ex.finished {
		return
	}

	if ex.fromPool && ex.vw.examplePool != nil {
		C.VW_ReturnExampleToPool(ex.vw.examplePool, ex.exHandle)
	} else {
		C.VW_FinishExample(ex.vw.handle, ex.exHandle)
	}
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

// GetAction returns the action for index
func (ex *Example) GetAction(i int) int {
	return int(C.VW_GetAction(ex.exHandle, C.size_t(i)))
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

// GetActions returns all actions as a slice
func (ex *Example) GetActions() []int {
	length := ex.GetActionScoreLength()
	if length < 0 {
		return nil
	}

	actions := make([]int, length)

	for i := 0; i < length; i++ {
		actions[i] = ex.GetAction(i)
	}

	return actions
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
func (ex *Example) GetConfidence() float32 {
	return float32(C.VW_GetConfidence(ex.exHandle))
}

// GetLoss returns loss for the example
func (ex *Example) GetLoss() float32 {
	return float32(C.VW_GetLoss(ex.exHandle))
}

// GetCBCostLength get cb cost length
func (ex *Example) GetCBCostLength() int {
	return int(C.VW_GetCBCostLength(ex.exHandle))
}

// GetCBCost get cb cost length
func (ex *Example) GetCBCost(i int) int {
	return int(C.VW_GetCBCost(ex.exHandle, C.size_t(i)))
}

// GetMultiClassPrediction returns prediction value for a multiclass prediction
func (ex *Example) GetMultiClassPrediction() int {
	return int(C.VW_GetMultiClassPrediction(ex.exHandle))
}

// GetScalarLength returns the length of scalar predictions in the example
func (ex *Example) GetScalarLength() int {
	return int(C.VW_GetScalarLength(ex.exHandle))
}

// GetScalar returns the length of scalar predictions in the example
func (ex *Example) GetScalar(i int) float32 {
	return float32(C.VW_GetScalar(ex.exHandle, C.size_t(i)))
}

// ExampleList a slice of examples
type ExampleList []*Example

// Finish all examples in the list
func (el ExampleList) Finish() {
	for _, e := range el {
		e.Finish()
	}
}
