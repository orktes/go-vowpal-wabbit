package vw

import (
	"fmt"
	"testing"
)

func TestNewVWInstance(t *testing.T) {
	vw, err := New("-q st --noconstant --quiet")
	if err != nil {
		t.Fatal(err)
	}
	defer vw.Finish()

	ex, err := vw.ReadExample("1 |s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	if err != nil {
		t.Fatal(err)
	}
	score := vw.Learn(ex)
	if score != 0 {
		t.Error("should start for zero")
	}

	score = vw.Learn(ex)
	if score == 0 {
		t.Error("should not be zero")
	}

	prediction := vw.Predict(ex)

	if prediction == 0 {
		t.Error("predict should be returned")
	}
}

func TestVWCopyModelData(t *testing.T) {
	vw, err := New("-q st --noconstant --quiet")
	if err != nil {
		t.Fatal(err)
	}
	defer vw.Finish()

	ex, err := vw.ReadExample("1 |s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	if err != nil {
		t.Fatal(err)
	}
	vw.Learn(ex)

	modelData := vw.CopyModelData()
	if len(modelData) == 0 {
		t.Error("model data should not be empty")
	}

	anotherModel, err := NewWithModelData("-q st --noconstant --quiet", modelData)
	if err != nil {
		t.Fatal(err)
	}
	defer anotherModel.Finish()

	ex, err = anotherModel.ReadExample("1 |s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	if err != nil {
		t.Fatal(err)
	}
	prediction := vw.Predict(ex)

	if prediction == 0 {
		t.Error("predict should be returned")
	}

}

func TestVWWithSeed(t *testing.T) {
	vw, err := New("-q st --noconstant --quiet")
	if err != nil {
		t.Fatal(err)
	}
	defer vw.Finish()

	ex, err := vw.ReadExample("1 |s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	if err != nil {
		t.Fatal(err)
	}
	vw.Learn(ex)

	anotherModel, err := NewWithSeed(vw, "")
	if err != nil {
		t.Fatal(err)
	}
	defer anotherModel.Finish()

	ex, err = anotherModel.ReadExample("1 |s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	if err != nil {
		t.Fatal(err)
	}
	prediction := vw.Predict(ex)

	if prediction == 0 {
		t.Error("predict should be returned")
	}
}

func ExampleVW() {
	vw, _ := New("-q st --noconstant --quiet")
	defer vw.Finish()

	// Learn an example
	ex, _ := vw.ReadExample("1 |s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	vw.Learn(ex)

	// Predict with features
	ex, _ = vw.ReadExample("|s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	res := vw.Predict(ex)

	fmt.Printf("Prediction: %f", res)
	// Output: Prediction: 0.855723
}
