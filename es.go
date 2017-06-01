// Package es provides an Elasticsearch query DSL.
package es

import (
	"encoding/json"
	"fmt"
	"strings"
)

// compress JSON.
func compress(s string) string {
	var v interface{}

	if err := json.Unmarshal([]byte(s), &v); err != nil {
		panic(err)
	}

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// When returns `children` only when `cond` is met.
func When(cond bool, children ...string) string {
	if cond {
		return join(children)
	}

	return ""
}

// Pretty JSON.
func Pretty(s string) string {
	var v interface{}

	if err := json.Unmarshal([]byte(s), &v); err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

// Query is the root of a query.
func Query(children ...string) string {
	return compress(fmt.Sprintf(`{
    "size": 0,
    %s
  }`, join(children)))
}

// DateHistogram applies a date_histogram.
func DateHistogram(interval string) string {
	return fmt.Sprintf(`
		"date_histogram": {
			"field": "timestamp",
			"interval": %q
		}
	`, interval)
}

// Filter applies the given filters.
func Filter(filters ...string) func(children ...string) string {
	return func(children ...string) string {
		return fmt.Sprintf(`
			"filter": {
				"bool": {
					"filter": [
						%s
					]
				}
			},
			%s
		`, join(filters), join(children))
	}
}

// Range for filtering.
func Range(gte, lte string) string {
	return fmt.Sprintf(`{
		"range": {
			"timestamp": {
				"gte": %q,
				"lte": %q
			}
		}
	}`, gte, lte)
}

// Term returns a term reference for filtering.
func Term(field, value string) string {
	return fmt.Sprintf(`{
		"term": {
			%q: %q
		}
	}`, field, value)
}

// Aggs of the given name.
func Aggs(name string, children ...string) string {
	return fmt.Sprintf(`
  "aggs": {
    %q: {
      %s
    }
  }`, name, join(children))
}

// Terms agg of the given field.
func Terms(field string, size int) string {
	return fmt.Sprintf(`
    "terms": {
      "field": %q,
      "size": %d
    }
  `, field, size)
}

// Sum agg of the given field.
func Sum(field string) string {
	return fmt.Sprintf(`
    "sum": {
      "field": %q
    }
  `, field)
}

// Avg agg of the given field.
func Avg(field string) string {
	return fmt.Sprintf(`
    "avg": {
      "field": %q
    }
  `, field)
}

// Min agg of the given field.
func Min(field string) string {
	return fmt.Sprintf(`
    "min": {
      "field": %q
    }
  `, field)
}

// Max agg of the given field.
func Max(field string) string {
	return fmt.Sprintf(`
    "max": {
      "field": %q
    }
  `, field)
}

// Stats agg of the given field.
func Stats(field string) string {
	return fmt.Sprintf(`
    "stats": {
      "field": %q
    }
  `, field)
}

// Percentiles agg of the given field, optionally
// specifying which `percents` to include.
func Percentiles(field string, percents ...float64) string {
	if len(percents) > 0 {
		return fmt.Sprintf(`
      "stats": {
        "field": %q,
        "percents": [%s]
      }
    `, field, joinFloats(percents))
	}

	return fmt.Sprintf(`
    "stats": {
      "field": %q
    }
  `, field)
}

// Histogram agg of the given field.
func Histogram(field string, options ...string) string {
	return fmt.Sprintf(`
    "histogram": {
      "field": %q,
      %s
    }
  `, field, join(options))
}

// Interval of `n`.
func Interval(n int) string {
	return fmt.Sprintf(`"interval": %d`, n)
}

// MinDocCount of `n`.
func MinDocCount(n int) string {
	return fmt.Sprintf(`"min_doc_count": %d`, n)
}

// Missing value of `n`.
func Missing(n int) string {
	return fmt.Sprintf(`"missing": %d`, n)
}

// ExtendedBounds of `min` / `max`.
func ExtendedBounds(min, max int) string {
	return fmt.Sprintf(`"extended_bounds": {
    "min": %d,
    "max": %d
  }`, min, max)
}

// Direction for sorting.
type Direction string

// Directions available.
const (
	Ascending  Direction = "asc"
	Descending           = "desc"
)

// Order `field` by `direction`.
func Order(field string, dir Direction) string {
	return fmt.Sprintf(`"order": {
    %q: %q
  }`, field, dir)
}

// JoinFloats returns floats joined by a comma.
func joinFloats(vals []float64) string {
	var s []string

	for _, v := range vals {
		s = append(s, fmt.Sprintf("%0.2f", v))
	}

	return strings.Join(s, ", ")
}

// Join returns strings joined by a comma.
func join(s []string) string {
	return strings.Join(clean(s), ",\n")
}

// Clean empty strings which may be present from When().
func clean(s []string) (vals []string) {
	for _, v := range s {
		if v != "" {
			vals = append(vals, v)
		}
	}
	return
}
