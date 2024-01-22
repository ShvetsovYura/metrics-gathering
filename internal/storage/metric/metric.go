package metric

import "strconv"

type Gauge float64
type Counter int64

func (g Gauge) ToString() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}
func (g Gauge) GetRawValue() float64 {
	return float64(g)
}

func (c Counter) ToString() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c Counter) GetRawValue() int64 {
	return int64(c)
}
