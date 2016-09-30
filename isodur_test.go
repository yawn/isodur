package isodur

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	assert := assert.New(t)
	e := errors.New

	var tt = []struct {
		in      string
		err     error
		out     float64
		reverse string
	}{
		{"1D", e(`unexpected initial period designator "1" at <input>:1:2`), 0, ""},
		{"P1X", e(`unexpected date period designator "X" at <input>:1:4`), 0, ""},
		{"PT1X", e(`unexpected time period designator "X" at <input>:1:5`), 0, ""},
		{"P1.5W2D", e(`unexpected next period "2" following a previous decimal period at <input>:1:7`), 0, ""},
		{"P279769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000D", e(`unexpected period "279769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" at <input>:1:311`), 0, ""},
		{"P365D", nil, 3.1536e+07, "P1Y"},
		{"P30D", nil, 2.592e+06, "P1M"},
		{"P7D", nil, 604800, "P1W"},
		{"PT24H", nil, 86400, "P1D"},
		{"PT60M", nil, 3600, "PT1H"},
		{"PT60S", nil, 60, "PT1M"},
		{"P1.75D", nil, 151200, "P1DT18H"},
		{"P0.75D", nil, 64800, "PT18H"},
		{"P3Y6M4DT12H30M5S", nil, (3*365*24*time.Hour + 6*30*24*time.Hour + 4*24*time.Hour + 12*time.Hour + 30*time.Minute + 5*time.Second).Seconds(), "P3Y6M4DT12H30M5S"},
		{"PT2H1,5M", nil, 7200 + 90, "PT2H1M30S"},
		{"PT1M1.23456S", nil, 61.23456, "PT1M1.23456S"},
		{"PT2.5S", nil, 2.5, "PT2.5S"},
	}

	for _, e := range tt {

		out, err := Parse(e.in)

		if err != nil || e.err != nil {
			assert.Equal(e.err, err, e.in)
		} else {
			assert.Equal(e.out, out.Seconds(), e.in)
			assert.Equal(e.reverse, out.String(), e.in)
		}

	}

}
