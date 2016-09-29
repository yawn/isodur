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

	const base = 10

	var (
		buf     = bytes.NewBufferString(p)
		days    int64
		dur     = int64(_d.Duration)
		hours   int64
		minutes int64
		months  int64
		seconds float64
		weeks   int64
		years   int64
	)

	// TODO: cleanup

	if dur >= int64(year) {

		years = dur / int64(year)
		dur = dur - (years * int64(year))

		buf.WriteString(strconv.FormatInt(years, base))
		buf.WriteString(y)

	}

	if dur >= int64(month) {

		months = dur / int64(month)
		dur = dur - (months * int64(month))

		buf.WriteString(strconv.FormatInt(months, base))
		buf.WriteString(m)

	}

	if dur >= int64(week) {

		weeks = dur / int64(week)
		dur = dur - (weeks * int64(week))

		buf.WriteString(strconv.FormatInt(weeks, base))
		buf.WriteString(w)

	}

	if dur >= int64(day) {

		days = dur / int64(day)
		dur = dur - (days * int64(day))

		buf.WriteString(strconv.FormatInt(days, base))
		buf.WriteString(d)

	}

	if dur >= int64(hour) {

		hours = dur / int64(hour)
		dur = dur - (hours * int64(hour))

		buf.WriteString(t)

		buf.WriteString(strconv.FormatInt(hours, base))
		buf.WriteString(h)

	}

	if dur >= int64(minute) {

		minutes = dur / int64(minute)
		dur = dur - (minutes * int64(minute))

		if hours == 0 {
			buf.WriteString(t)
		}

		buf.WriteString(strconv.FormatInt(minutes, base))
		buf.WriteString(m)

	}

	seconds = float64(dur) / float64(second)

	if seconds > 0 {

		if hours == 0 && minutes == 0 {
			buf.WriteString(t)
		}

		buf.WriteString(strconv.FormatFloat(seconds, 'f', -1, 64))
		buf.WriteString(s)

	}

	return buf.String()

}
