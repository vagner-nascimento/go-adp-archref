package rest

import (
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/data/httpdata"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"github.com/vagner-nascimento/go-enriching-adp/src/singleton"
)

type MerchantsClient struct {
	client *httpdata.HttpClient
}

var singMerchCli singleton.SingResource

func (mc *MerchantsClient) GetMerchant(id string) (merchant *appentity.Merchant, err error) {
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
