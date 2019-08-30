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

func TestReadDSJSON(t *testing.T) {
	input := `
	  {
		"_label_cost": -1,
		"_label_probability": 0.8166667,
		"_label_Action": 2,
		"_labelIndex": 1,
		"Version": "1",
		"EventId": "0074434d3a3a46529f65de8a59631939",
		"a": [
		  2,
		  1,
		  3
		],
		"c": {
		  "shared_ns": {
			"shared_feature": 0
		  },
		  "_multi": [
			{
			  "_tag": "tag",
			  "ns1": {
				"f1": 1,
				"f2": "strng"
			  },
			  "ns2": [
				{
				  "f3": "value1"
				},
				{
				  "ns3": {
					"f4": 0.994963765
				  }
				}
			  ]
			},
			{
			  "_tag": "tag",
			  "ns1": {
				"f1": 1,
				"f2": "strng"
			  }
			},
			{
			  "_tag": "tag",
			  "ns1": {
				"f1": 1,
				"f2": "strng"
			  }
			}
		  ]
		},
		"p": [
		  0.816666663,
		  0.183333333,
		  0.183333333
		],
		"VWState": {
		  "m": "096200c6c41e42bbb879c12830247637/0639c12bea464192828b250ffc389657"
		}
	  }
	`

	vw, _ := New("--dsjson --cb_adf --no_stdin")
	defer vw.Finish()

	examples, err := vw.ReadDecisionServiceJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	if len(examples) != 4 {
		t.Error("expecting 4 examples")
	}

}

func TestReadJSON(t *testing.T) {
	input := `
	{
		"_label": 1,
		"features": {
		  "13": 3.9656971e-02,
		  "24303": 2.2660980e-01,
		  "const": 0.01
		}
	}
	`

	vw, _ := New("--json --no_stdin")
	defer vw.Finish()

	examples, err := vw.ReadJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	if len(examples) != 1 {
		t.Error("should contain one example but contained", len(examples))
	}

	examples.Finish()
}

func TestReadJSONWithMultiLineLearn(t *testing.T) {
	input := `
	{"_labelIndex":1,"_label_Action":0,"_label_Cost":0,"_label_Probability":0.5,"_multi":[{"b_":"1","c_":"1","d_":"1"}, {"b_":"2","c_":"2","d_":"2"}]}
	`

	vw, _ := New("--json --cb_explore_adf --no_stdin")
	defer vw.Finish()

	examples, err := vw.ReadJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	vw.MultiLineLearn(examples)

	if examples[0].GetActionScore(0) != 0.5 {
		t.Error("should have been 0.5")
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
