package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

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
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	// Configure Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		APIKey: "VmhZUlA1VUJOSXlHQlNQZFlINzM6b0xYLTItN2tTbUdxZHJvbHR4d1dOdw==",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error().Msgf("Error creating Elasticsearch client: %s", err)
		return
	}

	indexName := "cooking_blog"

	logger.Info().Msgf("creating elasticsearch index: %s", indexName)

	// creating a new index (an index is a basic storage unit for a document).
	err = CreateNewIndex(es, indexName, logger)
	if err != nil {
		logger.Error().Msgf("Error creating index: %s", err)
		return
	}

	logger.Info().Msgf(" successfully created elasticsearch index: %s", indexName)

	// explicitly creating an index schema for the cooking_blog.
	err = PutNewIndexMapping(es, indexName, logger)
	if err != nil {
		logger.Error().Msgf("Error creating index mapping: %s", err)
		return
	}

	// adding documents to created elasticsearch index.
	recipes := []struct {
		Index  map[string]string `json:"index"`
		Recipe Recipe            `json:"-"`
	}{
		{Index: map[string]string{"_id": "1"}, Recipe: Recipe{
			Title:       "Perfect Pancakes: A Fluffy Breakfast Delight",
			Description: "Learn the secrets to making the fluffiest pancakes...",
			Author:      "Maria Rodriguez",
			Date:        "2023-05-01",
			Category:    "Breakfast",
			Tags:        []string{"pancakes", "breakfast", "easy recipes"},
			Rating:      4.8,
		}},
		{Index: map[string]string{"_id": "2"}, Recipe: Recipe{
			Title:       "Spicy Thai Green Curry: A Vegetarian Adventure",
			Description: "Dive into the flavors of Thailand with this vibrant green curry...",
			Author:      "Liam Chen",
			Date:        "2023-05-05",
			Category:    "Main Course",
			Tags:        []string{"thai", "vegetarian", "curry", "spicy"},
			Rating:      4.6,
		}},
	}
	var buf bytes.Buffer
	for _, r := range recipes {
		meta, _ := json.Marshal(r.Index)
		buf.Write(meta)
		buf.WriteByte('\n')

		// Serialize actual document
		data, _ := json.Marshal(r.Recipe)
		buf.Write(data)
		buf.WriteByte('\n')
	}

	err = AddDocuments(es, indexName, buf, logger)
	if err != nil {
		logger.Error().Msgf("Error adding documents: %s", err)
		return
	}

	logger.Info().Msgf("successfully added documents to elasticsearch index: %s", indexName)

	// Perform a search query
	// searchResp, err := es.Search(
	//     es.Search.WithContext(context.Background()),
	//     es.Search.WithIndex("documents"),
	//     es.Search.WithQuery("snow"),
	//     es.Search.WithTrackTotalHits(true),
	//     es.Search.WithPretty(),
	// )
	// Perform a search query
	// searchResp, err := es.Search(
	// 	es.Search.WithContext(context.Background()),
	// 	es.Search.WithIndex("documents"),
	// 	es.Search.WithQuery("snow"),
	// 	es.Search.WithTrackTotalHits(true),
	// 	es.Search.WithPretty(),
	// )

}

func AddDocuments(es *elasticsearch.Client, indexName string, data bytes.Buffer, logger zerolog.Logger) error {
	// Send bulk request to Elasticsearch
	res, err := es.Bulk(bytes.NewReader(data.Bytes()), es.Bulk.WithIndex(indexName), es.Bulk.WithRefresh("wait_for"))
	if err != nil {
		logger.Error().Msgf("Error sending bulk request: %s", err)
		return err
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

	logger.Info().Msgf("success creating index %s", res.Body)
	return nil
}
