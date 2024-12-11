package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"

	"github.com/gilbsgilbs/intervalst/interval"
)

func measure[T any](fn func() T) (T, time.Duration) {
	start := time.Now()
	result := fn()
	return result, time.Now().Sub(start)
}

func main() {
	ivMax := 1000000000
	maxIvDiff := 10000
	nbElems := 1000000
	nbMeasure := 10

	st := interval.NewMultiValueSearchTreeWithOptions[int, int](func(a, b int) int { return a - b }, interval.TreeWithIntervalPoint())
	for i := 0; i < nbElems; i++ {
		start := rand.Intn(ivMax) - maxIvDiff
		end := start + rand.Intn(maxIvDiff)
		itv := []int{start, end}
		st.Insert(itv[0], itv[1], i)
	}

	for tn, itvfunc := range []func() []int{
		func() []int { return []int{1, 10} },
		func() []int { return []int{ivMax - 1000, ivMax} },
		func() []int { r := max(0, rand.Intn(ivMax)-1000); return []int{r, r + 1000} },
		func() []int { return []int{rand.Intn(ivMax), rand.Intn(ivMax)} },
	} {
		var durationOld, durationNew time.Duration

		for i := 0; i < nbMeasure; i++ {
			var resOld, resNew []int
			var durationOld_, durationNew_ time.Duration

			itv := itvfunc()
			sort.Ints(itv)

			interval.SioImplNew = false
			resOld, durationOld_ = measure(func() []int {
				res, _ := st.AllIntersections(itv[0], itv[1])
				return res
			})
			durationOld += durationOld_

			interval.SioImplNew = true
			resNew, durationNew_ = measure(func() []int {
				res, _ := st.AllIntersections(itv[0], itv[1])
				return res
			})
			durationNew += durationNew_

			if !reflect.DeepEqual(resOld, resNew) {
				fmt.Printf("Got diff !!! old=%+v new=%+v", resOld, resNew)
				panic("")
			}
		}

		fmt.Printf("Test %d: %v => %v\n", tn, durationOld, durationNew)
	}
}
