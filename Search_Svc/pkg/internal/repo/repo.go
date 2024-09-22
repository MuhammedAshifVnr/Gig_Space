package repo

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"
	"github.com/elastic/go-elasticsearch/v8"
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
