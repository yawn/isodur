# isodur

[![Build Status](https://travis-ci.org/yawn/isodur.svg?branch=master)](https://travis-ci.org/yawn/isodur)

`isodir` is an ISO 8601 duration period parser and generator. It supports
fractions.

## Examples

* `Parse("PT24H")` yields a duration of 86400 seconds
* `Parse("PT60M")` yields a duration of 3600 seconds
* `Parse("P1.75D")` yields a duration of 151200 seconds
* `Parse("PT2H1,5M")` yields a duration of 7290 seconds
