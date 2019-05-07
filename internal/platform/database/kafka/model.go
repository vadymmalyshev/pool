package kafka

import (
	"errors"
)

// DB represents Kafka settings
type DB struct {
	Brokers     string `yaml:"brokers"`
	CaLocation  string `yaml:"ca_location"`
	Username    string `yaml:"username"`
	Pass        string `yaml:"password"`
	Topics      string `yaml:"topics"`
	GroupID     string `yaml:"group_id"`
	Retention   string `yaml:"retention"`
	DbName      string `yaml:"db_name"`
	Precision   string `yaml:"precision"`
	MiningPools string `yaml:"mining_pools"`
}

// Validate returns error if config is not valid
func (db *DB) Validate() error {
	if db.Brokers == "" {
		return errors.New("kafka brokers are required")
	}

	if db.CaLocation == "" {
		return errors.New("kafka certs are required")
	}

	if db.Username == "" {
		return errors.New("kafka user is required")
	}

	if db.Pass == "" {
		return errors.New("kafka password is required")
	}

	if db.Topics == "" {
		return errors.New("kafka topics are required")
	}

	if db.GroupID == "" {
		return errors.New("kafka group id is required")
	}

	if db.Retention == "" {
		return errors.New("kafka retention is required")
	}

	if db.DbName == "" {
		return errors.New("kafka db name is required")
	}

	if db.Precision == "" {
		return errors.New("kafka precision is required")
	}

	if db.MiningPools == "" {
		return errors.New("kafka mining pools are required")
	}

	return nil
}
