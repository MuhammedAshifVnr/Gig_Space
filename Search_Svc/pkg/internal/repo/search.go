package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/redis/go-redis/v9"
)

func (r *Repo) SearchGigs(ctx context.Context, query string, priceUpto float32, revisionsMin, deliveryDaysMax int32) ([]*model.Gig, error) {

	key := fmt.Sprintf("search:%s:%f:%d:%d", query, priceUpto, revisionsMin, deliveryDaysMax)
	result, err := r.GetCachedResponse(ctx, key)

	if err == nil && result != "" {
		log.Println("inside the redis")
		var gigs []*model.Gig
		if err := json.Unmarshal([]byte(result), &gigs); err == nil {
			return gigs, nil
		}
	}

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
			"filter": []interface{}{},
		},
	}

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

	queryBody, err := json.Marshal(map[string]interface{}{
		"query": boolQuery,
	})
	if err != nil {
		log.Println("Error marshaling query:", err)
		return nil, err
	}

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

	var gigs []*model.Gig
	for _, hit := range searchResults.Hits.Hits {
		var gig model.Gig
		if err := json.Unmarshal(hit.Source, &gig); err == nil {
			gigs = append(gigs, &gig)
		}
	}

	cacheData, _ := json.Marshal(gigs)
	r.CacheResponse(ctx, key, string(cacheData), time.Minute*10)

	return gigs, nil
}

func (r *Repo) GetCachedResponse(ctx context.Context, key string) (string, error) {
	result, err := r.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

func (r *Repo) CacheResponse(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.RDB.Set(ctx, key, value, expiration).Err()
}
