package kafka

import (
	"errors"
)

type Config struct {
	KafkaBrokers     string
	KafkaCaLocation  string
	KafkaUsername    string
	KafkaPass        string
	KafkaTopics      string
	KafkaGroupId     string
	KafkaRetention   string
	KafkaDbName      string
	KafkaPrecision   string
	KafkaMiningPools string
}

// Validate returns error if config is not valid
func (c Config) Validate() error {
	if c.KafkaBrokers == "" {
		return errors.New("kafka brokers are required")
	}

	if c.KafkaCaLocation == "" {
		return errors.New("kafka certs are required")
	}

	if c.KafkaUsername == "" {
		return errors.New("kafka user is required")
	}

	if c.KafkaPass == "" {
		return errors.New("kafka password is required")
	}

	if c.KafkaTopics == "" {
		return errors.New("kafka topics are required")
	}

	if c.KafkaGroupId == "" {
		return errors.New("kafka group id is required")
	}

	if c.KafkaRetention == "" {
		return errors.New("kafka retention is required")
	}

	if c.KafkaDbName == "" {
		return errors.New("kafka db name is required")
	}

	if c.KafkaPrecision == "" {
		return errors.New("kafka precision is required")
	}

	if c.KafkaMiningPools == "" {
		return errors.New("kafka mining pools are required")
	}

	return nil
}



