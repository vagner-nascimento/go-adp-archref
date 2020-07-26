package rest

import (
	"github.com/vagner-nascimento/go-adp-bridge/config"
	appentity "github.com/vagner-nascimento/go-adp-bridge/src/app/entity"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/httpdata"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-bridge/src/singleton"
)

type MerchantsClient struct {
	client *httpdata.HttpClient
}

var singMerchCli singleton.SingResource

func (mc *MerchantsClient) GetMerchant(id string) (merchant appentity.Merchant, err error) {
	var (
		status int
		data   []byte
	)

	status, data, err = mc.client.Get("", id)

	if err = handleResponse(status, err, data, "merchant"); err != nil {
		logger.Error("error on try to get merchant", err)
	} else {
		merchant, err = appentity.NewMerchant(data)
	}

	return
}

func GetMerchantsClient() *MerchantsClient {
	singMerchCli.Once.Do(func() {
		conf := config.Get().Integration.Rest.Merchants

		singMerchCli.Resource = &MerchantsClient{
			client: httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized),
		}
	})

	return singMerchCli.Resource.(*MerchantsClient)
}
