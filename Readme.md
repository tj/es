# ES

Package es provides an Elasticsearch query DSL.

## Example

If you don't mind crazy nesting:

```go
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
```

If you do mind crazy nesting:

```go
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
```

Both yielding:

```json
{
  "aggs": {
    "results": {
      "aggs": {
        "repos": {
          "aggs": {
            "labels": {
              "aggs": {
                "duration_sum": {
                  "sum": {
                    "field": "duration"
                  }
                }
              },
              "terms": {
                "field": "issue.labels.keyword",
                "size": 100
              }
            }
          },
          "terms": {
            "field": "repository.name.keyword",
            "size": 100
          }
        }
      },
      "filter": {
        "bool": {
          "filter": [
            {
              "term": {
                "user.login": "tj"
              }
            },
            {
              "range": {
                "timestamp": {
                  "gte": "now-7d",
                  "lte": "now"
                }
              }
            }
          ]
        }
      }
    }
  },
  "size": 0
}
```

---

[![GoDoc](https://godoc.org/github.com/tj/es?status.svg)](https://godoc.org/github.com/tj/es)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-experimental-orange.svg)

<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>
