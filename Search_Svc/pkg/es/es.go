package es

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

func InitElasticSearch() *elasticsearch.Client{
	cfg := elasticsearch.Config{Addresses: []string{viper.GetString("Url")}}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}
	return es
}
