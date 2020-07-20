package repository

import (
	"encoding/json"
	"errors"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/amqpdata"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-bridge/src/integration/rest"
	"net/http"
)

type accountRepository struct {
	topic           string
	merchantAccCli  *rest.MerchantAccountsClient
	merchantsCli    *rest.MerchantsClient
	affiliationsCli *rest.AffiliationsClient
}

func (repo *accountRepository) Save(account *app.Account) error {
	if bytes, err := json.Marshal(account); err == nil {
		return amqpdata.Publish(bytes, repo.topic)
	} else {
		return apperror.New("error on convert Account into bytes", err, nil)
	}
}

func (repo *accountRepository) GetMerchantAccount(accId string) (data []byte, err error) {
	status, data, gErr := repo.merchantAccCli.GetAccount(accId)

	msg := "error on try to get Merchant Account"

	if gErr != nil {
		err = apperror.New(msg, err, nil)

		return
	}

	if isHttpResponseFailed(status) {
		if status == http.StatusNotFound {
			err = handleNotfoundError("Merchant Account not found", data)
		} else {
			err = apperror.New(msg, nil, data)
		}

		data = nil
	}

	return
}

func (repo *accountRepository) GetMerchantAccounts(merchantId string) (data []byte, err error) {
	status, data, gErr := repo.merchantAccCli.GetAccountList(map[string]string{"merchant_id": merchantId})
	msg := "error on try to get Merchant Accounts"

	if gErr != nil {
		err = apperror.New(msg, gErr, nil)
		return
	}

	if isHttpResponseFailed(status) {
		if status == http.StatusNotFound {
			err = handleNotfoundError("Merchant accounts not found", data)
		} else {
			err = apperror.New(msg, nil, data)
		}

		data = nil
	}

	return
}

func (repo *accountRepository) GetMerchant(merchantId string) (data []byte, err error) {
	// TODO: handle http question into client
	status, data, gErr := repo.merchantsCli.GetMerchant(merchantId)
	msg := "error on try to get Merchant"

	if gErr != nil {
		logger.Error(msg, gErr)
		err = errors.New(msg)
		return
	}

	if isHttpResponseFailed(status) {
		if status == http.StatusNotFound {
			err = handleNotfoundError("Merchant accounts not found", data)
		} else {
			err = apperror.New(msg, nil, data)
			logger.Error(msg, err)
		}

		data = nil
	}

	return
}

func (repo *accountRepository) GetAffiliation(affId string) (data []byte, err error) {
	status, data, gErr := repo.affiliationsCli.GetAffiliation(affId)
	msg := "error on try to get Affiliation"

	if gErr != nil {
		logger.Error(msg, gErr)
		err = errors.New(msg)
		return
	}

	if isHttpResponseFailed(status) {
		if status == http.StatusNotFound {
			err = handleNotfoundError("Affiliation not found", data)
		} else {
			err = apperror.New(msg, nil, data)
			logger.Error(msg, err)
		}

		data = nil
	}

	return
}

func NewAccountRepository() *accountRepository {
	return &accountRepository{
		topic:           config.Get().Integration.Amqp.Pubs.CrmAccount.Topic,
		merchantAccCli:  rest.GetMerchantAccClient(),
		merchantsCli:    rest.GetMerchantsClient(),
		affiliationsCli: rest.GetAffiliationsClient(),
	}
}
