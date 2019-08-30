package main

import (
	"encoding/json"
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

func serializeToJSON(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func main() {
	vw, err := wabbit.New("--cb_explore 4 --cover 3 --json --no_stdin")
	if err != nil {
		panic(err)
	}
	defer vw.Finish()

	contextValues := []string{"a", "b", "c", "d"}

	wins := 0
	tries := 0

	for {
		contextValue := randomValue(contextValues)

		predictExamples, err := vw.ReadJSON(serializeToJSON(map[string]interface{}{
			"f": map[string]interface{}{
				"some_feature": contextValue,
			},
		}))
		if err != nil {
			panic(err)
		}

		predictExample := predictExamples[0]

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

		trainExamples, err := vw.ReadJSON(serializeToJSON(map[string]interface{}{
			"_label_Action":      (selected + 1),
			"_label_Cost":        cost,
			"_label_Probability": scores[selected],
			"f": map[string]interface{}{
				"some_feature": contextValue,
			},
		}))
		if err != nil {
			panic(err)
		}
		trainExample := trainExamples[0]

		vw.Learn(trainExample)
		trainExample.Finish()

		tries++
		if (tries < 500 && tries%50 == 0) || tries%500 == 0 {
			fmt.Printf("try %d: selected right arm %d out of %d (%f)\n", tries, wins, tries, (float64(wins)/float64(tries))*100)
		}

		if tries > 10000 {
			break
		}
	}
}
