package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/redis/go-redis/v9"
)

type Repo struct {
	esClient *elasticsearch.Client
	RDB      *redis.Client
}

func NewSearchRepository(es *elasticsearch.Client, rdb *redis.Client) *Repo {
	return &Repo{
		esClient: es,
		RDB:      rdb,
	}
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

func (r *Repo) UpdateGig(data model.Gig, index string) error {
	docID, err := r.GetDocId(uint(data.Id), index)
	if err != nil {
		return err
	}
	updatedDataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	updateReq := esapi.IndexRequest{
		Index:      index,
		DocumentID: docID,
		Body:       bytes.NewReader(updatedDataJSON),
		Refresh:    "true",
	}

	updateRes, err := updateReq.Do(context.Background(), r.esClient)
	if err != nil {
		return err
	}
	defer updateRes.Body.Close()

	if updateRes.IsError() {
		log.Printf("Error updating gig ID %d: %s", data.Id, updateRes.String())
		return fmt.Errorf("failed to update gig: %s", updateRes.String())
	}

	return nil
}

func (r *Repo) GetDocId(gigID uint, index string) (string, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": gigID,
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	searchRes, err := r.esClient.Search(
		r.esClient.Search.WithContext(context.Background()),
		r.esClient.Search.WithIndex(index),
		r.esClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return "", err
	}
	defer searchRes.Body.Close()

	var searchResult map[string]interface{}
	if err := json.NewDecoder(searchRes.Body).Decode(&searchResult); err != nil {
		return "", err
	}

	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return "", fmt.Errorf("no document found for gig ID %d", gigID)
	}
	docID := hits[0].(map[string]interface{})["_id"].(string)
	return docID, nil
}

func (r *Repo) DeleteDocument(docID, index string) error {
	deleteReq := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docID,
	}

	deleteRes, err := deleteReq.Do(context.Background(), r.esClient)
	if err != nil {
		return err
	}
	defer deleteRes.Body.Close()

	if deleteRes.IsError() {
		log.Printf("Error deleting gig ID %s: %s", docID, deleteRes.String())
		return fmt.Errorf("failed to delete gig: %s", deleteRes.String())
	}

	return nil
}
