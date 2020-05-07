package repository

import (
	"encoding/json"
	"errors"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/applicationerror"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/rabbitmq"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/rest"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"net/http"
	"time"
)

type accountRepository struct {
	topic          string
	merchantAccCli *rest.Client
	merchantsCli   *rest.Client
}

func (repo *accountRepository) Save(account *app.Account) (err error) {
	if bytes, err := json.Marshal(account); err == nil {
		err = rabbitmq.Publish(bytes, repo.topic)
	} else {
		err = applicationerror.New("error on convert account's interface into bytes", err, nil)
	}

	return err
}

func (repo *accountRepository) GetMerchantAccounts(merchantId string) (data []byte, err error) {
	status, data, err := repo.merchantAccCli.GetMany("", map[string]string{"merchant_id": merchantId})
	msg := "error on try to get Merchant Accounts"

	if err != nil {
		logger.Error(msg, err)
		err = errors.New(msg)
	}

	if isHttpResponseFailed(status) {
		if status == http.StatusNotFound {
			err = handleNotfoundError("Merchant accounts not found", data)
		} else {
			err = applicationerror.New(msg, nil, data)
			logger.Error(msg, err)
		}
		data = nil
	}

	return data, err
}

func (repo *accountRepository) GetMerchant(merchantId string) (data []byte, err error) {
	status, data, err := repo.merchantsCli.Get("", merchantId)
	msg := "error on try to get Merchant"

	if err != nil {
		logger.Error(msg, err)
		err = errors.New(msg)
	}

	if isHttpResponseFailed(status) {
		if status == http.StatusNotFound {
			err = handleNotfoundError("Merchant accounts not found", data)
		} else {
			err = applicationerror.New(msg, nil, data)
			logger.Error(msg, err)
		}
		data = nil
	}

	return data, err
}

func NewAccountRepository() *accountRepository {
	intConf := config.Get().Integration
	mAccCliConf := intConf.Rest.MerchantAccounts
	mCliConf := intConf.Rest.Merchants

	return &accountRepository{
		topic: intConf.Amqp.Pubs.CrmAccount.Topic,
		merchantAccCli: rest.NewClient(
			mAccCliConf.BaseUrl,
			mAccCliConf.TimeOut*time.Second,
			mAccCliConf.RejectUnauthorized,
		),
		merchantsCli: rest.NewClient(
			mCliConf.BaseUrl,
			mCliConf.TimeOut*time.Second,
			mCliConf.RejectUnauthorized,
		),
	}
}
