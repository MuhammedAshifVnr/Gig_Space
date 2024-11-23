package es

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitElasticSearch() *elasticsearch.Client {
	cfg := elasticsearch.Config{Addresses: []string{viper.GetString("URL")}}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}
	return es
}

func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}
