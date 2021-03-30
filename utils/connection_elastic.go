package utils

import (
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

func NewElasticSearchClient(host, user, pass string) *elasticsearch.Client {
	addresses := []string{}
	for _, address := range strings.Split(host, ",") {
		addresses = append(addresses, "http://"+address)
	}
	log.Print(addresses)
	var cfg elasticsearch.Config
	if user == "" && pass == "" {
		cfg = elasticsearch.Config{
			Addresses: addresses,
		}
	} else {
		cfg = elasticsearch.Config{
			Addresses: addresses,
			Username:  user,
			Password:  pass,
		}
	}
	Es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Connect to elasticsearch failed: %s", err)
	}
	_, err = Es.Info()
	if err != nil {
		log.Fatalf("Get elasticsearch info failed: %s", err)
	}
	log.Print("Elasticsearch client created")
	return Es
}
