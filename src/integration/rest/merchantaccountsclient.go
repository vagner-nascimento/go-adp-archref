package rest

import (
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/httpdata"
)

type MerchantAccountsClient struct {
	client *httpdata.HttpClient
}

func (mc *MerchantAccountsClient) GetAccount(id string) (int, []byte, error) {
	return mc.client.Get("", id)
}

func (mc *MerchantAccountsClient) GetAccountList(params map[string]string) (int, []byte, error) {
	return mc.client.GetMany("", params)
}

var singMerchAccCli singletonResource

func GetMerchantAccClient() *MerchantAccountsClient {
	singMerchAccCli.once.Do(func() {
		conf := config.Get().Integration.Rest.MerchantAccounts

		singMerchAccCli.resource = &MerchantAccountsClient{
			client: httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized),
		}
	})

	return singMerchAccCli.resource.(*MerchantAccountsClient)
}
