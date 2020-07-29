package repository

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/apperror"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/data/amqpdata"
	restintegration "github.com/vagner-nascimento/go-enriching-adp/src/integration/rest"
)

type accountRepository struct {
	topic           string
	publisher       *amqpdata.AmqpPublisher
	merchantAccCli  *restintegration.MerchantAccountsClient
	merchantsCli    *restintegration.MerchantsClient
	affiliationsCli *restintegration.AffiliationsClient
}

func (repo *accountRepository) Save(account *appentity.Account) error {
	var (
		bytes       []byte
		err         error
		isPublished bool
	)

	if bytes, err = json.Marshal(account); err == nil {
		isPublished, err = repo.publisher.Publish(bytes)

		if err == nil && !isPublished {
			err = apperror.New("account not saved", nil, nil)
		}
	}

	return err
}

func (repo *accountRepository) GetMerchantAccount(accId string) (appentity.MerchantAccount, error) {
	return repo.merchantAccCli.GetMerchantAccount(accId)
}

func (repo *accountRepository) GetMerchantAccounts(merchantId string) ([]appentity.MerchantAccount, error) {
	return repo.merchantAccCli.GetMerchantAccounts(map[string]string{"merchant_id": merchantId})
}

func (repo *accountRepository) GetMerchant(merchantId string) (*appentity.Merchant, error) {
	return repo.merchantsCli.GetMerchant(merchantId)
}

func (repo *accountRepository) GetAffiliation(affId string) (appentity.Affiliation, error) {
	return repo.affiliationsCli.GetAffiliation(affId)
}

func NewAccountRepository() *accountRepository {
	return &accountRepository{
		publisher:       amqpdata.NewAmqpPublisher(config.Get().Integration.Amqp.Pubs.CrmAccount.Topic),
		merchantAccCli:  restintegration.GetMerchantAccClient(),
		merchantsCli:    restintegration.GetMerchantsClient(),
		affiliationsCli: restintegration.GetAffiliationsClient(),
	}
}
