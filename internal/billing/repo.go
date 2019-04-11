package billing

import (
	"fmt"
	"log"
	"time"

	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	"git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type BillingRepositorer interface {
	SaveWallet(wallet models.Wallet) error
	BulkUpdateWallets(wallets []models.Wallet) error
	BulkUpdateWorkers(workers []models.Worker) error
	CreateWorkerIfNotExists(worker string) (*models.Worker, error)
	CreateWalletIfNotExists(wallet string, coinName string) (*models.Wallet, error)
	CreateCoinIfNotExists(coin string) (*models.Coin, error)
	DeleteWallet(name string)
	FindWalletByName(name string) (models.Wallet, error)
	SaveWorker(worker models.Worker) error
	DeleteWorker(name string)
	FindWorkerByName(name string) (models.Worker, error)
	GetWalletEarning(wallet string, date string) (WalletEarning, error)
	SaveWorkerStatistic(workerStatistic models.BillingWorkerStatistic, wallet string, worker string) (*models.Worker,
		error)
	AllWorkerStatistic(workerIDs []int, date string) ([]models.BillingWorkerStatistic, error)
	SaveWorkerMoney(moneyStatistic models.BillingWorkerMoney) error
	BulkUpdateWorkersShares(fees []models.WorkerFee) error
	BulkUpdateWorkersFeeIfNotExist(fees []models.WorkerFee) error
	FindPaidWorkersFee(formattedDate string, workerID int) error
	FindWorkersFeeByDate(date string) ([]models.WorkerFee, error)
}

type BillingRepository struct {
	client *gorm.DB
}

func NewBillingRepository() *BillingRepository {
	return &BillingRepository{GetBillingRepositoryClient()}
}

func GetBillingRepositoryClient() *gorm.DB {
	db, err := postgres.Connect(config.DB)

	//err = models.Migrate(db) // testing

	if err != nil {
		log.Panic("failed to init billing db :", err.Error())
	}
	return db
}

func (r *BillingRepository) SaveWallet(wallet models.Wallet) error {
	return r.client.FirstOrCreate(&wallet, wallet).Error
}

func (r *BillingRepository) BulkUpdateWallets(wallets []models.Wallet) error {
	fmt.Println("Started wallets updates ", time.Now())
	for _, v := range wallets {
		if err := r.client.FirstOrCreate(&v, v).Error; err != nil {
			return err
		}
	}
	fmt.Println("Finished wallets updates", time.Now())
	return nil
}

func (r *BillingRepository) BulkUpdateWorkers(workers []models.Worker) error {
	fmt.Println("Started workers updates ", time.Now())
	for _, v := range workers {
		if err := r.client.FirstOrCreate(&v, v).Error; err != nil {
			return err
		}
	}
	fmt.Println("Finished workers updates", time.Now())
	return nil
}

func (r *BillingRepository) CreateWorkerIfNotExists(worker string) (*models.Worker, error) {
	var billingWorker models.Worker
	res := r.client.FirstOrCreate(&billingWorker, models.Worker{Name: worker})
	if res.Error != nil {
		return nil, res.Error
	}
	return &billingWorker, nil
}

func (r *BillingRepository) CreateWalletIfNotExists(wallet string, coinName string) (*models.Wallet, error) {
	var billingWallet models.Wallet
	coin, err := r.CreateCoinIfNotExists(coinName)
	if err != nil {
		return nil, err
	}
	w := models.Wallet{Address: wallet, CoinID: coin.ID}
	res := r.client.FirstOrCreate(&billingWallet, w)
	if res.Error != nil {
		return nil, res.Error
	}

	return &billingWallet, nil
}

func (r *BillingRepository) CreateCoinIfNotExists(coin string) (*models.Coin, error) {
	var billingCoin models.Coin
	res := r.client.FirstOrCreate(&billingCoin, models.Coin{Name: coin})
	if res.Error != nil {
		return nil, res.Error
	}
	return &billingCoin, nil
}

func (r *BillingRepository) DeleteWallet(name string) {
	r.client.Unscoped().Where("name LIKE ?", name).Delete(models.Wallet{})
	r.client.Commit()
}

func (r *BillingRepository) FindWalletByName(name string) (models.Wallet, error) {
	var billingWallet models.Wallet
	notFound := r.client.Where("name = ?", name).First(&billingWallet).RecordNotFound()

	var err error

	if notFound {
		err = fmt.Errorf("can find wallet with name %s", name)
	}
	return billingWallet, err
}

func (r *BillingRepository) SaveWorker(worker models.Worker) error {
	return r.client.FirstOrCreate(&worker, worker).Error
}

func (r *BillingRepository) DeleteWorker(name string) {
	r.client.Unscoped().Where("name LIKE ?", name).Delete(models.Worker{})
	r.client.Commit()
}

func (r *BillingRepository) FindWorkerByName(name string) (models.Worker, error) {
	var billingWorker models.Worker
	notFound := r.client.Where("name = ?", name).First(&billingWorker).RecordNotFound()

	var err error

	if notFound {
		err = fmt.Errorf("can find wallet with name %s", name)
	}
	return billingWorker, err
}

func (r *BillingRepository) GetWalletEarning(wallet string, date string) (WalletEarning, error) {
	var result Earning

	sql := fmt.Sprintf(` select w.address, sum(hashrate) as hashrate, sum(usd) as usd, sum(cny) as cny, sum(btc) as btc, sum(commission_usd) as commission
	from billing_money m
	join billing_statistic b
	on b.worker_id = m.worker_id
	join wallets w
	on w.id = b.wallet_id
	where to_char((m.created_at - INTERVAL '1 DAY'), 'DD-Mon-YYYY') = '%s'
    and w.address = '%s'
    group by w.address`, date, wallet)

	err := r.client.Raw(sql).Row().Scan(&result.Address, &result.Hashrate, &result.USD, &result.CNY, &result.BTC, &result.Commission_USD)
	if apierrors.HandleError(err) {
		return WalletEarning{}, err
	}
	result.Date = date
	return WalletEarning{200, result}, nil
}

func (r *BillingRepository) SaveWorkerStatistic(workerStatistic models.BillingWorkerStatistic, wallet string,
	worker string) (*models.Worker, error) {
	dbWorker, _ := r.CreateWorkerIfNotExists(worker)
	dbWallet, _ := r.CreateWalletIfNotExists(wallet, "ETH")
	workerStatistic.Worker = *dbWorker
	workerStatistic.Wallet = *dbWallet
	r.client.NewRecord(workerStatistic)
	return dbWorker, r.client.Create(&workerStatistic).Error
}

// AllWorkerStatistic returns AllWorkerStatistic.
func (r *BillingRepository) AllWorkerStatistic(workerIDs []int, date string) ([]models.BillingWorkerStatistic, error) {
	var stats []models.BillingWorkerStatistic

	sql := fmt.Sprintf(`SELECT bs.wallet_id, bs.worker_id, bs.valid_shares, wf.coin, wf.wallet_address, wr.name
		FROM wallets AS w
			   JOIN billing_statistic AS bs ON w.id = bs.wallet_id
			   JOIN worker_fees AS wf ON w.address = wf.wallet_address
               JOIN workers AS wr ON wr.id = bs.worker_id
		WHERE wf.billing_date = '%s'`, date)

	rows, err := r.client.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	for rows.Next() {
		var stat models.BillingWorkerStatistic
		if err := rows.Scan(&stat.WalletID, &stat.WorkerID, &stat.ValidShares,
			&stat.Wallet.Coin.Name, &stat.Wallet.Address, &stat.Worker.Name); err != nil {
			logrus.Error(err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *BillingRepository) SaveWorkerMoney(moneyStatistic models.BillingWorkerMoney) error {
	r.client.NewRecord(moneyStatistic)
	return r.client.Create(&moneyStatistic).Error
}

// FindPaidWorkersFee finds workers fee.
// formattedDate - 02.01.2006
func (r *BillingRepository) FindPaidWorkersFee(formattedDate string, workerID int) error {
	var fee models.WorkerFee
	notFound := r.client.Where(models.WorkerFee{Date: formattedDate, WorkerID: workerID, Paid: true}).
		First(&fee).RecordNotFound()

	if notFound {
		return fmt.Errorf("can't find any worker statistic")
	}

	return nil
}

// BulkUpdateWorkersShares update workers shares if not exist.
func (r *BillingRepository) BulkUpdateWorkersShares(fees []models.WorkerFee) error {
	fmt.Println("Started charges updates ", time.Now())
	tx := r.client.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			logrus.Errorf("error transaction: %v", err)
		}
	}()
	for _, v := range fees {
		err := tx.
			Model(&v).
			Update(models.WorkerFee{Shares: v.Shares, Paid: true}).
			Error
		if err != nil {
			return err
		}
	}
	err := tx.Commit().Error
	fmt.Println("Finished charges updates", time.Now())

	return err
}

// BulkUpdateWorkersFeeIfNotExists upserts workers fee if not exist.
func (r *BillingRepository) BulkUpdateWorkersFeeIfNotExist(fees []models.WorkerFee) error {
	fmt.Println("Started charges updates ", time.Now())
	tx := r.client.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			logrus.Errorf("error transaction: %v", err)
		}
	}()
	for _, v := range fees {
		if err := r.FindPaidWorkersFee(v.Date, v.WorkerID); err != nil {
			err := tx.
				Where(models.WorkerFee{Date: v.Date, WorkerID: v.WorkerID, Paid: false}).
				Assign(models.WorkerFee{Amount: v.Amount, Coin: v.Coin, WalletAddr: v.WalletAddr}).
				FirstOrCreate(&v).
				Error
			if err != nil {
				return err
			}
		}
	}
	err := tx.Commit().Error
	fmt.Println("Finished charges updates", time.Now())

	return err
}

// FindWorkersFeeByDate finds fees for worker by date
func (r *BillingRepository) FindWorkersFeeByDate(date string) ([]models.WorkerFee, error) {
	var fees []models.WorkerFee

	notFound := r.client.Find(&fees, models.WorkerFee{Date: date}).RecordNotFound()

	if notFound {
		return nil, fmt.Errorf("can't find any worker statistic")
	}

	return fees, nil
}
