package es_test

import (
	"fmt"

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
