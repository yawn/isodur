package isodur

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	sc "text/scanner"
	"time"
)

const (
	day    = 24 * time.Hour
	hour   = time.Hour
	minute = time.Minute
	month  = 30 * day
	second = time.Second
	week   = 7 * day
	year   = 365 * day
)

const (
	comma = ","
	d     = "D"
	dot   = "."
	h     = "H"
	p     = "P"
	t     = "T"
	m     = "M"
	s     = "S"
	w     = "W"
	y     = "Y"
)

// Duration extends a time.Duration with ISO 8601 formatting and parsing
// capabilities.
type Duration struct {
	time.Duration
}

// Parse parses the given ISO 8601 duration string
func Parse(str string) (*Duration, error) {

	var (
		duration int64
		fraction float64
		scanner  sc.Scanner
		timeMode = false
		token    rune
		total    time.Duration
	)

	scanner.Init(bytes.NewBufferString(strings.Replace(str, comma, dot, -1)))
	scanner.Mode = sc.ScanChars

	for token != sc.EOF {

		token = scanner.Scan()
		text := scanner.TokenText()

		if text == "" {
			token = sc.EOF
		} else if scanner.Pos().Offset == 1 {

			if text != p {
				return nil, fmt.Errorf("unexpected initial period designator %q at %s", text, scanner.Pos())
			}

			scanner.Mode = sc.ScanFloats

		} else if text == t {
			timeMode = true
		} else if scanner.Mode == sc.ScanChars {

			var unit time.Duration

			if !timeMode {

				switch text {
				case y:
					unit = year
				case m:
					unit = month
				case w:
					unit = week
				case d:
					unit = day
				default:
					return nil, fmt.Errorf("unexpected date period designator %q at %s", text, scanner.Pos())
				}

			} else {

				switch text {
				case h:
					unit = hour
				case m:
					unit = minute
				case s:
					unit = second
				default:
					return nil, fmt.Errorf("unexpected time period designator %q at %s", text, scanner.Pos())
				}

			}

			total = total + time.Duration(duration)*unit + time.Duration(fraction*float64(unit))

			if scanner.Peek() != sc.Char {
				scanner.Mode = sc.ScanFloats
			}

		} else {

			if fraction > 0 {
				return nil, fmt.Errorf("unexpected next period %q following a previous decimal period at %s", text, scanner.Pos())
			}

			n, err := strconv.ParseFloat(text, 10)

			if err != nil {
				return nil, fmt.Errorf("unexpected period %q at %s", text, scanner.Pos())
			}

			d, f := math.Modf(n)

			duration = int64(d)
			fraction = f

			scanner.Mode = sc.ScanChars

		}

	}

	return &Duration{total}, nil

}

// String formats the duration to ISO 8601.
func (_d *Duration) String() string {

	var (
		buf      = bytes.NewBufferString(p)
		dur      = int64(_d.Duration)
		slots    = make([]float64, 7)
		timeMode = false
		tokens   = []string{y, m, w, d, h, m, s}
		windows  = []time.Duration{year, month, week, day, hour, minute, second}
	)

	for i := range slots {

		n := float64(windows[i])
		d := int64(n)

		if dur >= d {

			var units float64

			if i < 6 { // try to avoid fractionals
				units = float64(dur / d)
				dur = dur - (int64(units) * d)
			} else {
				units = float64(dur) / n
				dur = int64(float64(dur) - units*n)
			}

			if i > 3 && !timeMode && units > 0 {
				timeMode = true
				buf.WriteString(t)
			}

			buf.WriteString(strconv.FormatFloat(float64(units), 'f', -1, 64))
			buf.WriteString(tokens[i])

		}

	}

	return buf.String()

}
