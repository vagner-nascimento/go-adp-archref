package repository

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/amqpdata"
	"github.com/vagner-nascimento/go-adp-bridge/src/integration/rest"
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

func (repo *accountRepository) GetMerchantAccount(accId string) (app.MerchantAccount, error) {
	return repo.merchantAccCli.GetMerchantAccount(accId)
}

func (repo *accountRepository) GetMerchantAccounts(merchantId string) ([]app.MerchantAccount, error) {
	return repo.merchantAccCli.GetMerchantAccounts(map[string]string{"merchant_id": merchantId})
}

func (repo *accountRepository) GetMerchant(merchantId string) (app.Merchant, error) {
	return repo.merchantsCli.GetMerchant(merchantId)
}

func (repo *accountRepository) GetAffiliation(affId string) (app.Affiliation, error) {
	return repo.affiliationsCli.GetAffiliation(affId)
}

func NewAccountRepository() *accountRepository {
	return &accountRepository{
		topic:           config.Get().Integration.Amqp.Pubs.CrmAccount.Topic,
		merchantAccCli:  rest.GetMerchantAccClient(),
		merchantsCli:    rest.GetMerchantsClient(),
		affiliationsCli: rest.GetAffiliationsClient(),
	}
}
