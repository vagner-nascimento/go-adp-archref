package rest

import (
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/data/httpdata"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"github.com/vagner-nascimento/go-enriching-adp/src/singleton"
)

type MerchantAccountsClient struct {
	client *httpdata.HttpClient
}

func (mc *MerchantAccountsClient) GetMerchantAccount(id string) (mAcc appentity.MerchantAccount, err error) {
	var (
		status int
		data   []byte
	)

	status, data, err = mc.client.Get("", id)

	if err = handleResponse(status, err, data, "merchant account"); err != nil {
		logger.Error("error on try to get merchant account", err)
	} else {
		mAcc, err = appentity.NewMerchantAccount(data)
	}

	return
}

func (mc *MerchantAccountsClient) GetMerchantAccounts(params map[string]string) (mAccs []appentity.MerchantAccount, err error) {
	var (
		status int
		data   []byte
	)

	status, data, err = mc.client.GetMany("", params)

	if err = handleResponse(status, err, data, "merchant account"); err != nil {
		logger.Error("error on try to get merchant account", err)
	} else {
		mAccs, err = appentity.NewMerchantAccounts(data)
	}

	return
}

var singMerchAccCli singleton.SingResource

func GetMerchantAccClient() *MerchantAccountsClient {
	singMerchAccCli.Once.Do(func() {
		conf := config.Get().Integration.Rest.MerchantAccounts

		singMerchAccCli.Resource = &MerchantAccountsClient{
			client: httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized),
		}
	})

	return singMerchAccCli.Resource.(*MerchantAccountsClient)
}
