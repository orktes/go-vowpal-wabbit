package main

import (
	"fmt"
	"math/rand"

	wabbit "github.com/orktes/go-vowpal-wabbit"
)

func indexOf(val string, vals []string) int {
	for i, v := range vals {
		if val == v {
			return i
		}
	}

	return -1
}

func randomValue(vals []string) string {
	return vals[rand.Intn(len(vals))]
}

func argMax(vals []float32) int {
	maxIndx := 0
	maxVal := float32(0)

	for i, val := range vals {
		if val > maxVal {
			maxIndx = i
		}
	}
	return maxIndx
}

func main() {
	vw, err := wabbit.New("--cb_explore 4")
	if err != nil {
		panic(err)
	}
	defer vw.Finish()

	contextValues := []string{"a", "b", "c", "d"}

	wins := 0
	tries := 0

	for {
		contextValue := randomValue(contextValues)

		predictExample, err := vw.ReadExample(" | " + contextValue)
		if err != nil {
			panic(err)
		}

		vw.Predict(predictExample)

		scores := predictExample.GetActionScores()
		fmt.Printf("scores: %+v\n", scores)
		selected := argMax(scores)

		cost := 0.0
		if selected != indexOf(contextValue, contextValues) {
			cost = 0.5
		} else {
			cost = 0.0
			wins++
		}
		predictExample.Finish()

		example := fmt.Sprintf(" %d:%f:1 | %s", (selected + 1), cost, contextValue)
		println(example)
		trainExample, err := vw.ReadExample(example)
		if err != nil {
			panic(err)
		}

		vw.Learn(trainExample)
		trainExample.Finish()

		tries++

		if tries%1000 == 0 {
			fmt.Printf("%f\n", float64(wins)/float64(tries))
		}

	}
}
