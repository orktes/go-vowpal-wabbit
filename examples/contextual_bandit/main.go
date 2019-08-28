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
			maxVal = val
		}
	}
	return maxIndx
}

func weightedRandom(vals []float32) int {
	sum := float32(0)
	for _, val := range vals {
		sum += val
	}

	r := rand.Float32() * sum

	for i, val := range vals {
		r -= val
		if r <= 0 {
			return i
		}
	}

	return 0
}

func main() {
	vw, err := wabbit.New("--cb_explore 4 --cover 3")
	if err != nil {
		panic(err)
	}
	defer vw.Finish()

	contextValues := []string{"a", "b", "c", "d"}

	wins := 0
	tries := 0

	for {
		contextValue := randomValue(contextValues)

		predictExample, err := vw.ReadExample(fmt.Sprintf(" | %s", contextValue))
		if err != nil {
			panic(err)
		}

		vw.Predict(predictExample)

		scores := predictExample.GetActionScores()
		selected := weightedRandom(scores)

		cost := 0.0
		if selected != indexOf(contextValue, contextValues) {
			cost = 1.0
		} else {
			cost = 0.0
			wins++
		}
		predictExample.Finish()

		example := fmt.Sprintf(" %d:%f:%f | %s", (selected + 1), cost, scores[selected], contextValue)
		trainExample, err := vw.ReadExample(example)
		if err != nil {
			panic(err)
		}

		vw.Learn(trainExample)
		trainExample.Finish()

		tries++
		if tries%500 == 0 {
			fmt.Printf("try %d: selected right arm %d out of %d (%f)\n", tries, wins, tries, (float64(wins)/float64(tries))*100)
		}

		if tries > 10000 {
			break
		}
	}
}
