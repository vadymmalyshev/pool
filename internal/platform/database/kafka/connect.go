package kafka

import (
	client "github.com/confluentinc/confluent-kafka-go/kafka"
)

// Connect returns kafka consumer
func (db *DB) Connect() (*client.Consumer, error) {
	if err := db.Validate(); err != nil {
		return nil, err
	}
	config := &client.ConfigMap{
		"api.version.request":  "true",
		"metadata.broker.list": db.Brokers,
		"security.protocol":    "sasl_ssl",
		"sasl.mechanisms":      "PLAIN",
		"ssl.ca.location":      db.CaLocation,
		"sasl.username":        db.Username,
		"sasl.password":        db.Pass,
		"group.id":             db.GroupID,
		//"go.events.channel.enable":        true,
		//"go.application.rebalance.enable": true,
		"default.topic.config": client.ConfigMap{"auto.offset.reset": "earliest"},
	}
	cons, err := client.NewConsumer(config)
	return cons, err
}
