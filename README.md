<a id="markdown-rolling" name="rolling"></a>
# rolling
[![GoDoc](https://godoc.org/github.com/smallnest/rolling?status.svg)](https://godoc.org/github.com/smallnest/rolling)
[![Build Status](https://travis-ci.com/smallnest/rolling.png?branch=master)](https://travis-ci.com/smallnest/rolling)
[![codecov.io](https://codecov.io/github/smallnest/rolling/coverage.svg?branch=master)](https://codecov.io/github/smallnest/rolling?branch=master)

**A rolling/sliding window implementation for Google-golang**

<!-- TOC -->

- [rolling](#rolling)
    - [Usage](#usage)
        - [Point Window](#point-window)
        - [Time Window](#time-window)
    - [Aggregating Windows](#aggregating-windows)
            - [Custom Aggregations](#custom-aggregations)
    - [Contributors](#contributors)
    - [License](#license)

<!-- /TOC -->

<a id="markdown-usage" name="usage"></a>
## Usage

<a id="markdown-point-window" name="point-window"></a>
### Point Window

```golang
var p = rollingwindow.NewPointPolicy[float64](rolling.NewWindow[float64](5))

for x := 0; x < 5; x++{
  p.Append(x)
}
p.Reduce(func(w Window[float64]) float64 {
  fmt.Println(w) // [ [0] [1] [2] [3] [4] ]
  return 0
})
w.Append(5)
p.Reduce(func(w Window[float64]) float64 {
  fmt.Println(w) // [ [5] [1] [2] [3] [4] ]
  return 0
})
w.Append(6)
p.Reduce(func(w Window[float64]) float64 {
  fmt.Println(w) // [ [5] [6] [2] [3] [4] ]
  return 0
})
```

The above creates a window that always contains 5 data points and then fills
it with the values 0 - 4. When the next value is appended it will overwrite
the first value. The window continuously overwrites the oldest value with the
latest to preserve the specified value count. This type of window is useful
for collecting data that have a known interval on which they are capture or
for tracking data where time is not a factor.

<a id="markdown-time-window" name="time-window"></a>
### Time Window

```golang
var p = rollingwindow.NewTimeWindow[float64](rolling.NewWindow[float64](3000), time.Millisecond)
var start = time.Now()
for range time.Tick(time.Millisecond) {
  if time.Since(start) > 3*time.Second {
    break
  }
  p.Append(1)
}
```

The above creates a time window that contains 3,000 buckets where each bucket
contains, at most, 1ms of recorded data. The subsequent loop populates each
bucket with exactly one measure (the value 1) and stops when the window is full.
As time progresses, the oldest values will be removed such that if the above
code performed a `time.Sleep(3*time.Second)` then the window would be empty
again.

The choice of bucket size depends on the frequency with which data are expected
to be recorded. On each increment of time equal to the given duration the window
will expire one bucket and purge the collected values. The smaller the bucket
duration then the less data are lost when a bucket expires.

This type of bucket is most useful for collecting real-time values such as
request rates, error rates, and latencies of operations.

<a id="markdown-aggregating-windows" name="aggregating-windows"></a>
## Aggregating Windows

Each window exposes a `Reduce(func(w Window) float64) float64` method that can
be used to aggregate the data stored within. The method takes in a function
that can compute the contents of the `Window` into a single value. For
convenience, this package provides some common reductions:

```golang
fmt.Println(p.Reduce(rolling.Count))
fmt.Println(p.Reduce(rolling.Avg))
fmt.Println(p.Reduce(rolling.Min))
fmt.Println(p.Reduce(rolling.Max))
fmt.Println(p.Reduce(rolling.Sum))
fmt.Println(p.Reduce(rolling.Percentile(99.9)))
fmt.Println(p.Reduce(rolling.FastPercentile(99.9)))
```

The `Avg`, `Min`, `Max`, and `Sum` each perform their expected
computation. 

<a id="markdown-custom-aggregations" name="custom-aggregations"></a>
#### Custom Aggregations

Any function that matches the form of `func(rolling.Window)float64` may be given
to the `Reduce` method of any window policy. The `Window` type is a named
version of `[][]float64`. Calling `len(window)` will return the number of
buckets. Each bucket is, itself, a slice of floats where `len(bucket)` is the
number of values measured within that bucket. Most aggregate will take the form
of:

```golang
func MyAggregate(w rolling.Window) float64 {
  for _, bucket := range w {
    for _, value := range bucket {
      // aggregate something
    }
  }
}
```


<a id="markdown-license" name="license"></a>
## License

based on [asecurityteam/rolling](https://github.com/asecurityteam/rolling).

Apache 2.0 licensed, see [LICENSE.txt](LICENSE.txt) file.
