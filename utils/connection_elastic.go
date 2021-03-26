package utils

import (
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

var Es elasticsearch.Client

func GetElasticSearchConnection(host, user, pass string) {
	cfg := elasticsearch.Config{
		Addresses: strings.Split(host, ","),
		Username:  user,
		Password:  pass,
	}
	Es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Connect to elasticsearch failed: %s", err)
	}
	_, err = Es.Info()
	if err != nil {
		log.Fatalf("Get elasticsearch info failed: %s", err)
	}
}
