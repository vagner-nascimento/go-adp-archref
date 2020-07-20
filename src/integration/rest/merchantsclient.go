package rest

import (
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/httpdata"
)

type MerchantsClient struct {
	client *httpdata.HttpClient
}

var singMerchCli singletonResource

func (mc *MerchantsClient) GetMerchant(id string) (int, []byte, error) {
	return mc.client.Get("", id)
}
func GetMerchantsClient() *MerchantsClient {
	singMerchCli.once.Do(func() {
		conf := config.Get().Integration.Rest.Merchants

		singMerchCli.resource = &MerchantsClient{
			client: httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized),
		}
	})

	return singMerchCli.resource.(*MerchantsClient)
}
