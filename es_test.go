package es_test

import (
	"fmt"
	"testing"

	"github.com/tj/assert"

	. "github.com/tj/es"
)

func Example() {
	query := Pretty(Query(
		Aggs("results",
			Filter(
				Term("user.login", "tj"),
				Range("now-7d", "now"),
			)(
				Aggs("repos",
					Terms("repository.name.keyword", 100),
					Aggs("labels",
						Terms("issue.labels.keyword", 100),
						Aggs("duration_sum", Sum("duration"))))))))

	fmt.Println(query)
	// Output:
	// {
	//   "aggs": {
	//     "results": {
	//       "aggs": {
	//         "repos": {
	//           "aggs": {
	//             "labels": {
	//               "aggs": {
	//                 "duration_sum": {
	//                   "sum": {
	//                     "field": "duration"
	//                   }
	//                 }
	//               },
	//               "terms": {
	//                 "field": "issue.labels.keyword",
	//                 "size": 100
	//               }
	//             }
	//           },
	//           "terms": {
	//             "field": "repository.name.keyword",
	//             "size": 100
	//           }
	//         }
	//       },
	//       "filter": {
	//         "bool": {
	//           "filter": [
	//             {
	//               "term": {
	//                 "user.login": "tj"
	//               }
	//             },
	//             {
	//               "range": {
	//                 "timestamp": {
	//                   "gte": "now-7d",
	//                   "lte": "now"
	//                 }
	//               }
	//             }
	//           ]
	//         }
	//       }
	//     }
	//   },
	//   "size": 0
	// }
}

func Example_expanded() {
	labels := Aggs("labels",
		Terms("issue.labels.keyword", 100),
		Aggs("duration_sum",
			Sum("duration")))

	repos := Aggs("repos",
		Terms("repository.name.keyword", 100),
		labels)

	filter := Filter(
		Term("user.login", "tj"),
		Range("now-7d", "now"))

	results := Aggs("results", filter(repos))

	query := Pretty(Query(results))

	fmt.Println(query)
	// Output:
	// {
	//   "aggs": {
	//     "results": {
	//       "aggs": {
	//         "repos": {
	//           "aggs": {
	//             "labels": {
	//               "aggs": {
	//                 "duration_sum": {
	//                   "sum": {
	//                     "field": "duration"
	//                   }
	//                 }
	//               },
	//               "terms": {
	//                 "field": "issue.labels.keyword",
	//                 "size": 100
	//               }
	//             }
	//           },
	//           "terms": {
	//             "field": "repository.name.keyword",
	//             "size": 100
	//           }
	//         }
	//       },
	//       "filter": {
	//         "bool": {
	//           "filter": [
	//             {
	//               "term": {
	//                 "user.login": "tj"
	//               }
	//             },
	//             {
	//               "range": {
	//                 "timestamp": {
	//                   "gte": "now-7d",
	//                   "lte": "now"
	//                 }
	//               }
	//             }
	//           ]
	//         }
	//       }
	//     }
	//   },
	//   "size": 0
	// }
}

func TestPercentiles(t *testing.T) {
	t.Run("without percents", func(t *testing.T) {
		s := Pretty(Query(Percentiles("load_time")))
		assert.Equal(t, Pretty(`{"size":0,"stats":{"field":"load_time"}}`), s)
	})

	t.Run("with percents", func(t *testing.T) {
		s := Pretty(Query(Percentiles("load_time", 95, 99, 99.9)))
		assert.Equal(t, Pretty(`{"size":0,"stats":{"field":"load_time","percents":[95,99,99.9]}}`), s)
	})
}

func TestHistogram(t *testing.T) {
	h := Histogram("load_time",
		Interval(50),
		MinDocCount(1),
		ExtendedBounds(0, 500))

	s := Query(h)
	assert.Equal(t, `{"histogram":{"extended_bounds":{"max":500,"min":0},"field":"load_time","interval":50,"min_doc_count":1},"size":0}`, s)
}
