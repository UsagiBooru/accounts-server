package server

import (
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

func NewElasticSearchClient(host, user, pass string) *elasticsearch.Client {
	addresses := []string{}
	for _, address := range strings.Split(host, ",") {
		addresses = append(addresses, "http://"+address)
	}
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
		Error("Connect to elasticsearch failed: " + err.Error())
		os.Exit(1)
	}
	_, err = Es.Info()
	if err != nil {
		Error("Get elasticsearch info failed: " + err.Error())
		os.Exit(1)
	}
	// Debug("Elasticsearch client created")
	return Es
}
