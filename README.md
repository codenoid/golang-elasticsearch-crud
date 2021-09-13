# Go Elasticsearch CRUD example

example

## Setup

```sh
# daemonize the elasticsearch server
# you can use tmux split and not to daemonize elasticsearch
sudo docker-compose up -d
```

## CRUD Example

You can directly run `main.go` file, use `go run .`

### Go Module Installation

`go get github.com/olivere/elastic/v7`

### Connection

```go
client, err := elastic.NewClient(
    elastic.SetURL("http://127.0.0.1:9200"),
    elastic.SetBasicAuth("elastic", "esPwd123"), // password come from docker-compose ENV
    elastic.SetSniff(false),
    elastic.SetHealthcheck(false))

if err != nil {
    panic(err)
}
```

### Create Document

```go
student := map[string]interface{}{
    "name":        "Agung",
    "study":       "law",
    "last_attend": int64(0),
}
studentID := "1"

// Insert Document
_, err = elasticClient1.Index().
    Index("student").
    Id(studentID).
    BodyJson(student).
    Do(context.TODO())
```

### Update Document

```go
// update document
student["last_attend"] = time.Now().Unix()

// Update Document by Id
_, err = elasticClient1.Update().
    Index("student").
    Id(studentID).
    Doc(student).
    Do(context.TODO())
```

### Read Documents

Query String reference : [open ref](https://www.elastic.co/guide/en/elasticsearch/reference/master/query-dsl-query-string-query.html#query-string-syntax)

```go
limit := 10
esQueryTest := []string{"*", "name:agung", `study:"law"`, `agung AND last_attend:[0 TO *]`, `name:"sunby"`}
for _, query := range esQueryTest {
    fmt.Println("query test:", query)
    searchResult, err := elasticClient1.Search().
        Index("student").
        Query(elastic.NewQueryStringQuery(query)).
        Size(limit).
        Sort("last_attend", false).
        Do(context.TODO())
    panicIfError(err)

    var ttyp map[string]interface{}
    for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
        if t, ok := item.(map[string]interface{}); ok {
            fmt.Println(t)
        }
    }
}
```

### Delete Document

```go
_, err = elasticClient1.Delete().
    Index("student").
    Id(studentID).
    Do(context.TODO())
```