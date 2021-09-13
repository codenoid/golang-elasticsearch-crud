package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/olivere/elastic/v7"
)

var elasticClient1 *elastic.Client

func main() {
	if newClient, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetBasicAuth("elastic", "esPwd123"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false)); err != nil {
		panic(err)
	} else {
		elasticClient1 = newClient
	}

	ctx := context.TODO()

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := elasticClient1.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := elasticClient1.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

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
		Do(ctx)
	panicIfError(err)
	flushESDB("student")

	// update document
	student["last_attend"] = time.Now().Unix()

	// Update Document by Id
	elasticClient1.Update().
		Index("student").
		Id(studentID).
		Doc(student).
		Do(ctx)
	panicIfError(err)
	flushESDB("student")

	elasticClient1.Refresh("student")

	fmt.Println("====query test====")

	limit := 10
	// https://www.elastic.co/guide/en/elasticsearch/reference/master/query-dsl-query-string-query.html#query-string-syntax
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

	fmt.Println("====query test end====")

	// panic(2)
	// need delay
	// _, err = elasticClient1.Delete().
	// 	Index("student").
	// 	Id(studentID).
	// 	Do(context.TODO())
	// panicIfError(err)
	// flushESDB("student")

}

func flushESDB(indexname string) error {
	_, err := elasticClient1.Flush().Index(indexname).Do(context.TODO())
	return err
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
