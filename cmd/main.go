package main

import (
	"context"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// Configure Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		APIKey: "VDdHUUk1VUI3aHFKbVN1ZWtBaUI6ZzhFVnFIQmNUM1NZQ3NtTWhqQUx1dw==",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Define bulk index data
	_ = `
{"index":{"_id":"9780553351927"}}
{"name":"Snow Crash","author":"Neal Stephenson","release_date":"1992-06-01","page_count":470, "_extract_binary_content":true, "_reduce_whitespace":true, "_run_ml_inference":true}
{"index":{"_id":"9780441017225"}}
{"name":"Revelation Space","author":"Alastair Reynolds","release_date":"2000-03-15","page_count":585, "_extract_binary_content":true, "_reduce_whitespace":true, "_run_ml_inference":true}
{"index":{"_id":"9780451524935"}}
{"name":"1984","author":"George Orwell","release_date":"1985-06-01","page_count":328, "_extract_binary_content":true, "_reduce_whitespace":true, "_run_ml_inference":true}
{"index":{"_id":"9781451673319"}}
{"name":"Fahrenheit 451","author":"Ray Bradbury","release_date":"1953-10-15","page_count":227, "_extract_binary_content":true, "_reduce_whitespace":true, "_run_ml_inference":true}
{"index":{"_id":"9780060850524"}}
{"name":"Brave New World","author":"Aldous Huxley","release_date":"1932-06-01","page_count":268, "_extract_binary_content":true, "_reduce_whitespace":true, "_run_ml_inference":true}
{"index":{"_id":"9780385490818"}}
{"name":"The Handmaid's Tale","author":"Margaret Atwood","release_date":"1985-06-01","page_count":311, "_extract_binary_content":true, "_reduce_whitespace":true, "_run_ml_inference":true}
`

	// Send bulk request to Elasticsearch
	// ingestResult, err := es.Bulk(
	// 	bytes.NewReader([]byte(bookData)),
	// 	es.Bulk.WithIndex("documents"),
	// 	es.Bulk.WithPipeline("ent-search-generic-ingestion"),
	// )

	// // Output result
	// fmt.Println(ingestResult, err)

	// performing basic elastic search impl:
	searchResp, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("documents"),
		es.Search.WithQuery("snow"),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	  )
	  
	  fmt.Println(searchResp, err)
}
