package client

import (
	kaf "git.tor.ph/hiveon/pool/internal/platform/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Connect(c kaf.Config) (*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"api.version.request":  "true",
		"metadata.broker.list": c.KafkaBrokers,
		"security.protocol":    "sasl_ssl",
		"sasl.mechanisms":      "PLAIN",
		"ssl.ca.location":      c.KafkaCaLocation,
		"sasl.username":        c.KafkaUsername,
		"sasl.password":        c.KafkaPass,
		"group.id":             c.KafkaGroupId,
		//"go.events.channel.enable":        true,
		//"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
	cons, err := kafka.NewConsumer(config)
	return cons, err
}
