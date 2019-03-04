package wallets

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"git.tor.ph/hiveon/pool/config"

	red "github.com/gomodule/redigo/redis"
)

// WorkerPulse represents worker's name and last time it was sent share to the pool.
// It is typically extracted from Redis where has been put by the Consumer from Kafka stream.
type WorkerPulse struct {
	Name     string
	LastSeen time.Time
	Online   bool
}

const tsToSec = 1000000000

// GetWorkersPulse returns list of WorkersPluse instances
func GetWorkersPulse(redisClient red.Conn, walletID string) (result map[string]WorkerPulse, err error) {
	data, err := redisClient.Do("HGETALL", "last-update:"+walletID)

	if err != nil {
		return nil, errors.Wrap(err, "can't get workers pulse data")
	}

	workers, ok := data.(map[string]string)
	if !ok {
		return nil, errors.New("workers pulse data has wrong format")
	}

	for k, v := range workers {
		ts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't parse timestamp data of worker \"%s\"", k))
		}

		timeStamp := time.Unix(ts/tsToSec, 0)
		isOnline := timeStamp.After(time.Now().Add(-time.Duration(config.WorkerOfflineMin) * time.Minute))

		result[k] = WorkerPulse{
			Name:     k,
			LastSeen: timeStamp,
			Online:   isOnline,
		}
	}

	return result, nil
}
