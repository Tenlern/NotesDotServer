package models

import (
	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func SetupES() error {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	_, err = es.Info()
	if err != nil {
		return err
	}

	ES = es
	return nil
}
