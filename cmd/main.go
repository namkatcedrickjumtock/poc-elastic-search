package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog"
)

// elasticsearch specifications -> https://github.com/elastic/elasticsearch-specification

type Recipe struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Date        string   `json:"date"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Rating      float64  `json:"rating"`
}

func main() {
	ctx := context.Background()

	// Configure Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		APIKey: "VmhZUlA1VUJOSXlHQlNQZFlINzM6b0xYLTItN2tTbUdxZHJvbHR4d1dOdw==",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	indexName := "cooking_blog"

	// creating a new index (an index is a basic storage unit for a document).
	// err = CreateNewIndex(es, indexName, logger)
	// if err != nil {
	// 	log.Fatalf("%s", err)
	// }

	// explicitly creating an index schema for the cooking_blog.
	// err = PutNewIndexMapping(es, indexName, logger)
	// if err != nil {
	// 	log.Fatalf("Error creating index mapping: %s", err)
	// }

	// 	data := `
	// {"index":{"_id":"1"}}
	// {"title":"Perfect Pancakes: A Fluffy Breakfast Delight","description":"Learn the secrets to making the fluffiest pancakes, so amazing you won't believe your tastebuds. This recipe uses buttermilk and a special folding technique.","author":"Maria Rodriguez","date":"2023-05-01","category":"Breakfast","tags":["pancakes","breakfast","brunch"],"rating":4.8,"prep_time":"15 min","cook_time":"10 min","servings":4}
	// {"index":{"_id":"2"}}
	// {"title":"Spicy Thai Green Curry: A Vegetarian Adventure","description":"Dive into the flavors of Thailand with this vibrant green curry. Packed with vegetables and aromatic herbs.","author":"Liam Chen","date":"2023-05-05","category":"Main Course","tags":["thai","vegetarian","curry"],"rating":4.6,"prep_time":"25 min","cook_time":"20 min","servings":4}
	// {"index":{"_id":"3"}}
	// {"title":"Classic Beef Stroganoff: A Creamy Comfort Food","description":"Indulge in this rich and creamy beef stroganoff. Tender strips of beef in a savory mushroom sauce.","author":"Emma Watson","date":"2023-05-10","category":"Main Course","tags":["beef","pasta","comfort food"],"rating":4.7,"prep_time":"20 min","cook_time":"25 min","servings":4}
	// {"index":{"_id":"4"}}
	// {"title":"Vegan Chocolate Avocado Mousse","description":"Discover the magic of avocado in this rich, vegan chocolate mousse. Creamy, indulgent, and secretly healthy.","author":"Alex Green","date":"2023-05-15","category":"Dessert","tags":["vegan","chocolate","avocado"],"rating":4.5,"prep_time":"10 min","cook_time":"None","servings":2}
	// {"index":{"_id":"5"}}
	// {"title":"Crispy Oven-Fried Chicken","description":"Get that perfect crunch without the deep fryer! This oven-fried chicken recipe delivers crispy, juicy results.","author":"Maria Rodriguez","date":"2023-05-20","category":"Main Course","tags":["chicken","oven-fried","healthy"],"rating":4.9,"prep_time":"15 min","cook_time":"40 min","servings":4}
	// {"index":{"_id":"6"}}
	// {"title":"Homemade Margherita Pizza","description":"Simple and delicious Margherita pizza with a perfect balance of fresh tomato sauce, mozzarella, and basil.","author":"Luca Romano","date":"2023-06-02","category":"Main Course","tags":["pizza","Italian","cheese"],"rating":4.8,"prep_time":"20 min","cook_time":"12 min","servings":4}
	// {"index":{"_id":"7"}}
	// {"title":"Authentic Ramen Noodles","description":"Rich and flavorful Japanese ramen with homemade broth, tender pork, and perfectly cooked noodles.","author":"Kenji Tanaka","date":"2023-06-10","category":"Main Course","tags":["ramen","japanese","broth"],"rating":4.7,"prep_time":"30 min","cook_time":"4 hours","servings":4}
	// {"index":{"_id":"8"}}
	// {"title":"Strawberry Shortcake Bliss","description":"A classic strawberry shortcake recipe with fluffy biscuits, fresh strawberries, and sweet whipped cream.","author":"Jessica Carter","date":"2023-07-01","category":"Dessert","tags":["strawberries","shortcake","whipped cream"],"rating":4.9,"prep_time":"15 min","cook_time":"20 min","servings":6}
	// {"index":{"_id":"9"}}
	// {"title":"Gourmet Mac and Cheese","description":"A rich and creamy mac and cheese recipe with three types of cheese and a crispy breadcrumb topping.","author":"Michael Johnson","date":"2023-07-10","category":"Main Course","tags":["cheese","mac and cheese","comfort food"],"rating":4.6,"prep_time":"15 min","cook_time":"30 min","servings":4}
	// {"index":{"_id":"10"}}
	// {"title":"Classic French Croissants","description":"Buttery, flaky croissants made from scratch using traditional French pastry techniques.","author":"Sophie Dubois","date":"2023-07-20","category":"Breakfast","tags":["croissants","French","pastry"],"rating":4.9,"prep_time":"12 hours","cook_time":"20 min","servings":6}
	// {"index":{"_id":"30"}}
	// {"title":"Traditional Italian Tiramisu","description":"A rich and creamy tiramisu made with espresso-soaked ladyfingers, mascarpone cheese, and cocoa powder.","author":"Luca Romano","date":"2023-12-25","category":"Dessert","tags":["tiramisu","coffee","Italian"],"rating":4.9,"prep_time":"20 min","cook_time":"None","servings":8}
	// `

	// 	// Send bulk request to Elasticsearch
	// 	err = AddDocuments(es, indexName, data, logger)
	// 	if err != nil {
	// 		log.Fatalf("Error adding documents: %s", err)
	// 	}

	// Define the search query for basic full text search
	// https://www.elastic.co/guide/en/elasticsearch/reference/8.17/full-text-filter-tutorial.html#full-text-filter-tutorial-match-query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"description": map[string]string{
					"query": "Indulge in this rich",
				},
			},
		},
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error marshaling json query: %s", err)
	}

	// Perform a search query
	searchResp, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(indexName),
		// es.Search.WithQuery(searchTerm),
		es.Search.WithBody(bytes.NewReader(jsonQuery)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("error searching failed, why=%v, status=%d", err, searchResp.StatusCode)
	}

	fmt.Println(searchResp)
}

func AddDocuments(es *elasticsearch.Client, indexName string, data string, logger zerolog.Logger) error {
	// Send bulk request to Elasticsearch
	res, err := es.Bulk(
		bytes.NewReader([]byte(data)),
		es.Bulk.WithIndex(indexName),
		es.Bulk.WithPipeline("ent-search-generic-ingestion"),
	)
	if err != nil {
		return fmt.Errorf("error sending bulk request: %s", err)
	}

	logger.Info().Msgf("Successfully sent bulk request, status=%s", res.Status())
	return nil
}

// PutNewIndexMapping updates the mapping for an index
func PutNewIndexMapping(es *elasticsearch.Client, indexName string, logger zerolog.Logger) error {
	mapping := map[string]interface{}{
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type":     "text",
				"analyzer": "standard",
				"fields": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type":         "keyword",
						"ignore_above": 256,
					},
				},
			},
			"description": map[string]interface{}{
				"type": "text",
				"fields": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
			"author": map[string]interface{}{
				"type": "text",
				"fields": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
			"date": map[string]interface{}{
				"type":   "date",
				"format": "yyyy-MM-dd",
			},
			"category": map[string]interface{}{
				"type": "text",
				"fields": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
			"tags": map[string]interface{}{
				"type": "text",
				"fields": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
			"rating": map[string]interface{}{
				"type": "float",
			},
		},
	}

	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("error encoding mapping to JSON: %w", err)
	}

	res, err := es.Indices.PutMapping([]string{"cooking_blog"}, bytes.NewReader(mappingJSON))
	if err != nil {
		return fmt.Errorf("error sending request to update mapping: %w", err)
	}

	logger.Info().Msgf("successfully created elasticsearch index mapping for: %s: status=%s", indexName, res.Status())

	defer res.Body.Close()
	return nil
}

// CreateNewIndex creates a new index in Elasticsearch
func CreateNewIndex(es *elasticsearch.Client, indexName string, logger zerolog.Logger) error {
	res, err := es.Indices.Create(indexName)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request for index name %s status=%d", indexName, res.StatusCode)
	}

	logger.Info().Msgf("success creating index %s: status=%s", indexName, res.Status())

	return nil
}
