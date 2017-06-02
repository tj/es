package es_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/tj/assert"

	. "github.com/tj/es"
)

func Example() {
	query := Pretty(Query(
		Aggs(
			Agg("results",
				Filter(
					Term("user.login", "tj"),
					Range("now-7d", "now"),
				)(
					Aggs(
						Agg("repos",
							Terms("repository.name.keyword", 100),
							Aggs(
								Agg("labels",
									Terms("issue.labels.keyword", 100),
									Aggs(
										Agg("duration_sum", Sum("duration"))))))))))))

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
	labels := Aggs(
		Agg("labels",
			Terms("issue.labels.keyword", 100),
			Aggs(
				Agg("duration_sum",
					Sum("duration")))))

	repos := Aggs(
		Agg("repos",
			Terms("repository.name.keyword", 100),
			labels))

	filter := Filter(
		Term("user.login", "tj"),
		Range("now-7d", "now"))

	results := Aggs(
		Agg("results",
			filter(repos)))

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

func ExampleWhen() {
	period := "month"

	query := Pretty(Query(
		Aggs(
			Agg("results",
				Filter(
					Term("user.login", "tj"),
					When(period == "week", Range("now-7d", "now")),
					When(period == "month", Range("now-1M", "now")),
				)(
					Aggs(
						Agg("repos",
							Terms("repository.name.keyword", 100),
							Aggs(
								Agg("labels",
									Terms("issue.labels.keyword", 100),
									Aggs(
										Agg("duration_sum", Sum("duration"))))))))))))

	fmt.Println(query)
	// Output:
	// 	{
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
	//                   "gte": "now-1M",
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

func ExampleQuery() {
	s := Pretty(Query(
		Aggs(
			Agg("foo", Sum("count")),
			Agg("bar", Sum("count")),
			Agg("baz", Sum("count")),
		),
	))

	fmt.Println(s)
	// Output:
	// 	{
	//   "aggs": {
	//     "bar": {
	//       "sum": {
	//         "field": "count"
	//       }
	//     },
	//     "baz": {
	//       "sum": {
	//         "field": "count"
	//       }
	//     },
	//     "foo": {
	//       "sum": {
	//         "field": "count"
	//       }
	//     }
	//   },
	//   "size": 0
	// }

}

func TestPercentiles(t *testing.T) {
	t.Run("without percents", func(t *testing.T) {
		s := Query(Percentiles("load_time"))
		assert.Equal(t, `{"size":0,"stats":{"field":"load_time"}}`, s)
	})

	t.Run("with percents", func(t *testing.T) {
		s := Query(Percentiles("load_time", 95, 99, 99.9))
		assert.Equal(t, `{"size":0,"stats":{"field":"load_time","percents":[95,99,99.9]}}`, s)
	})
}

func TestDateHistogram(t *testing.T) {
	h := DateHistogram("timestamp",
		Interval("30m"),
		TimeZone("-08:00"))

	s := Query(h)
	assert.Equal(t, `{"date_histogram":{"field":"timestamp","interval":"30m","time_zone":"-08:00"},"size":0}`, s)
}

func TestHistogram(t *testing.T) {
	h := Histogram("load_time",
		Interval(50),
		MinDocCount(1),
		ExtendedBounds(0, 500),
		Order("something", Ascending))

	s := Query(h)
	assert.Equal(t, `{"histogram":{"extended_bounds":{"max":500,"min":0},"field":"load_time","interval":50,"min_doc_count":1,"order":{"something":"asc"}},"size":0}`, s)
}

func TestTimeZone(t *testing.T) {
	assert.Equal(t, `"time_zone": "-08:00"`, TimeZone("-08:00"), "offset string")
	assert.Equal(t, `"time_zone": "+05:45"`, TimeZone("Asia/Kathmandu"), "location")

	os.Setenv("TZ", "Asia/Kathmandu")
	assert.Equal(t, `"time_zone": "+05:45"`, TimeZone(), "default location")
}
