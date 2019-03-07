package billing

import (
	"fmt"
	"git.tor.ph/hiveon/pool/models"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	. "git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type BillingRepository struct {
	client *gorm.DB
}

func NewBillingRepository() *BillingRepository {
	return &BillingRepository{GetBillingRepositoryClient()}
}

func GetBillingRepositoryClient() *gorm.DB {
	db, err := postgres.Connect(config.DB)

	err = models.Migrate(db) // testing

	if err != nil {
		log.Panic("failed to init billing db :", err.Error())
	}
	return db
}

func (r *BillingRepository) SaveWallet(wallet Wallet) error {
	return r.client.FirstOrCreate(&wallet, wallet).Error
}

func (r *BillingRepository) BulkUpdateWallets(wallets []Wallet) error {
	fmt.Println("Started wallets updates ", time.Now())
	for _, v := range wallets {
		if err := r.client.FirstOrCreate(&v, v).Error; err != nil {
			return err
		}
	}
	fmt.Println("Finished wallets updates", time.Now())
	return nil
}

func (r *BillingRepository) BulkUpdateWorkers(workers []Worker) error {
	fmt.Println("Started workers updates ", time.Now())
	for _, v := range workers {
		if err := r.client.FirstOrCreate(&v, v).Error; err != nil {
			return err
		}
	}
	fmt.Println("Finished workers updates", time.Now())
	return nil
}

func (r *BillingRepository) CreateWorkerIfNotExists(worker string) (*Worker, error) {
	var billingWorker Worker
	res := r.client.FirstOrCreate(&billingWorker, Worker{Name: worker})
	if res.Error != nil {
		return nil, res.Error
	}
	return &billingWorker, nil
}

func (r *BillingRepository) CreateWalletIfNotExists(wallet string, coinName string) (*Wallet, error) {
	var billingWallet Wallet
	coin, err :=  r.CreateCoinIfNotExists(coinName)
	if err != nil {
		return nil, err
	}
	w := Wallet{Address: wallet, CoinID: coin.ID}
	res := r.client.FirstOrCreate(&billingWallet, w)
	if res.Error != nil {
		return nil, res.Error
	}

	return &billingWallet, nil
}

func (r *BillingRepository) CreateCoinIfNotExists(coin string) (*Coin, error) {
	var billingCoin Coin
	res := r.client.FirstOrCreate(&billingCoin, Coin{Name: coin})
	if res.Error != nil {
		return nil, res.Error
	}
	return &billingCoin, nil
}

func (r *BillingRepository) DeleteWallet(name string) {
	r.client.Unscoped().Where("name LIKE ?", name).Delete(Wallet{})
	r.client.Commit()
}

func (r *BillingRepository) FindWalletByName(name string) (Wallet, error) {
	var billingWallet Wallet
	notFound := r.client.Where("name = ?", name).First(&billingWallet).RecordNotFound()

	var err error

	if notFound {
		err = fmt.Errorf("Can find wallet with name %d", name)
	}
	return billingWallet, err
}

func (r *BillingRepository) SaveWorker(worker Worker) error {
	return r.client.FirstOrCreate(&worker, worker).Error
}

func (r *BillingRepository) DeleteWorker(name string) {
	r.client.Unscoped().Where("name LIKE ?", name).Delete(Worker{})
	r.client.Commit()
}

func (r *BillingRepository) FindWorkerByName(name string) (Worker, error) {
	var billingWorker Worker
	notFound := r.client.Where("name = ?", name).First(&billingWorker).RecordNotFound()

	var err error

	if notFound {
		err = fmt.Errorf("Can find wallet with name %d", name)
	}
	return billingWorker, err
}

func (r *BillingRepository) SaveWorkerStatistic(workerStatistic BillingWorkerStatistic, wallet string, worker string) (*Worker, error) {
	dbWorker, _ := r.CreateWorkerIfNotExists(worker)
	dbWallet, _ := r.CreateWalletIfNotExists(wallet, "ETH") // TODO:
	workerStatistic.Worker = *dbWorker
	workerStatistic.Wallet = *dbWallet
	r.client.NewRecord(workerStatistic)
	return dbWorker, r.client.Create(&workerStatistic).Error
}

func (r *BillingRepository) SaveWorkerMoney(moneyStatistic BillingWorkerMoney) error {
	r.client.NewRecord(moneyStatistic)
	return r.client.Create(&moneyStatistic).Error
}



