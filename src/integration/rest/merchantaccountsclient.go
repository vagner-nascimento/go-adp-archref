package restintegration

import (
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/data/httpdata"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"sync"
)

type MerchantAccountsClient struct {
	client  *httpdata.HttpClient
	connect sync.Once
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

var singletonMerchAccCli = &MerchantAccountsClient{}

func GetMerchantAccClient() *MerchantAccountsClient {
	singletonMerchAccCli.connect.Do(func() {
		conf := config.Get().Integration.Rest.MerchantAccounts
		singletonMerchAccCli.client = httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized)
	})

	return singletonMerchAccCli
}
