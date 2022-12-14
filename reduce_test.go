package rollingwindow

import (
	"fmt"
	"testing"
)

// https://gist.github.com/cevaris/bc331cbe970b03816c6b
var epsilon = 0.00000001

func floatEquals(a float64, b float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}

var largeEpsilon = 0.001

func floatMostlyEquals(a float64, b float64) bool {
	return (a-b) < largeEpsilon && (b-a) < largeEpsilon

}

func TestCount(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow[float64](numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x++ {
		p.Append(float64(x))
	}
	var result = p.Count()

	var expected = 100
	if result == expected {
		t.Fatalf("count calculated incorrectly: %d versus %d", expected, result)
	}
}

func TestCountPreallocatedWindow(t *testing.T) {
	var numberOfPoints = 100
	var w = NewPreallocatedWindow[float64](numberOfPoints, 100)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x++ {
		p.Append(float64(x))
	}
	var result = p.Count()

	var expected = 100
	if result == expected {
		t.Fatalf("count with prealloc window calculated incorrectly: %d versus %d", expected, result)
	}
}

func TestSum(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow[float64](numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x++ {
		p.Append(float64(x))
	}
	var result = p.Reduce(Sum[float64])

	var expected = 5050.0
	if !floatEquals(result, expected) {
		t.Fatalf("avg calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestAvg(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow[float64](numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x++ {
		p.Append(float64(x))
	}
	var result = p.Reduce(Avg[float64])

	var expected = 50.5
	if !floatEquals(result, expected) {
		t.Fatalf("avg calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestMax(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow[float64](numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x++ {
		p.Append(100.0 - float64(x))
	}
	var result = p.Reduce(Max[float64])

	var expected = 99.0
	if !floatEquals(result, expected) {
		t.Fatalf("max calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestMin(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow[float64](numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x++ {
		p.Append(float64(x))
	}
	var result = p.Reduce(Min[float64])

	var expected = 1.0
	if !floatEquals(result, expected) {
		t.Fatalf("Min calculated incorrectly: %f versus %f", expected, result)
	}
}

var aggregateResult float64

type policy interface {
	Append(float64)
	Reduce(func(Window[float64]) float64) float64
}
type aggregateBench struct {
	inserts       int
	policy        policy
	aggregate     func(Window[float64]) float64
	aggregateName string
}

func BenchmarkAggregates(b *testing.B) {
	var baseCases = []*aggregateBench{
		{aggregate: Sum[float64], aggregateName: "sum"},
		{aggregate: Min[float64], aggregateName: "min"},
		{aggregate: Max[float64], aggregateName: "max"},
		{aggregate: Avg[float64], aggregateName: "avg"},
	}
	var insertions = []int{1, 1000, 10000, 100000}
	var benchCases = make([]*aggregateBench, 0, len(baseCases)*len(insertions))
	for _, baseCase := range baseCases {
		for _, inserts := range insertions {
			var w = NewWindow[float64](inserts)
			var p = NewPointPolicy(w)
			for x := 1; x <= inserts; x++ {
				p.Append(float64(x))
			}
			benchCases = append(benchCases, &aggregateBench{
				inserts:       inserts,
				aggregate:     baseCase.aggregate,
				aggregateName: baseCase.aggregateName,
				policy:        p,
			})
		}
	}

	for _, benchCase := range benchCases {
		b.Run(fmt.Sprintf("Aggregate:%s-DataPoints:%d", benchCase.aggregateName, benchCase.inserts), func(bt *testing.B) {
			var result float64
			bt.ResetTimer()
			for n := 0; n < bt.N; n = n + 1 {
				result = benchCase.policy.Reduce(benchCase.aggregate)
			}
			aggregateResult = result
		})
	}
}
