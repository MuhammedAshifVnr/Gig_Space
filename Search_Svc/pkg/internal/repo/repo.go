package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Repo struct {
	esClient *elasticsearch.Client
}

func NewSearchRepository(es *elasticsearch.Client) *Repo {
	return &Repo{esClient: es}
}

func (r *Repo) IndexGig(data model.Gig, index string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = r.esClient.Index(index, strings.NewReader(string(jsonData)))
	if err != nil {
		log.Printf("Failed to index gig: %s", err)
		return err
	}

	return nil
}

func (r *Repo) SearchGigs(ctx context.Context, query string, priceUpto float32, revisionsMin, deliveryDaysMax int32) ([]*model.Gig, error) {
	boolQuery := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []interface{}{
				map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":     query,
						"fields":    []string{"title", "description", "category"},
						"fuzziness": "AUTO",
					},
				},
			},
			"filter": []interface{}{},  // Ensure filter is an array
		},
	}

	// Adding filter conditions
	if priceUpto > 0 {
		boolQuery["bool"].(map[string]interface{})["filter"] = append(boolQuery["bool"].(map[string]interface{})["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"lte": priceUpto,
				},
			},
		})
	}

	if revisionsMin > 0 {
		boolQuery["bool"].(map[string]interface{})["filter"] = append(boolQuery["bool"].(map[string]interface{})["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"revisions": map[string]interface{}{
					"gte": revisionsMin,
				},
			},
		})
	}

	if deliveryDaysMax > 0 {
		boolQuery["bool"].(map[string]interface{})["filter"] = append(boolQuery["bool"].(map[string]interface{})["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"delivery_days": map[string]interface{}{
					"lte": deliveryDaysMax,
				},
			},
		})
	}

	// Marshal the query
	queryBody, err := json.Marshal(map[string]interface{}{
		"query": boolQuery,  // Ensure "query" wraps the bool query
	})
	if err != nil {
		log.Println("Error marshaling query:", err)
		return nil, err
	}

	// Execute the search query
	req := esapi.SearchRequest{
		Index: []string{"gig"},
		Body:  bytes.NewReader(queryBody),
	}

	res, err := req.Do(ctx, r.esClient)
	if err != nil {
		log.Println("Search error:", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response: %s", res.Status())
		return nil, nil
	}

	// Parse the search results
	var searchResults struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResults); err != nil {
		log.Println("Error parsing search response:", err)
		return nil, err
	}

	// Unmarshal the results into a list of gigs
	var gigs []*model.Gig
	for _, hit := range searchResults.Hits.Hits {
		var gig model.Gig
		if err := json.Unmarshal(hit.Source, &gig); err == nil {
			gigs = append(gigs, &gig)
		}
	}

	return gigs, nil
}
