package restintegration

import (
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/data/httpdata"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"sync"
)

type MerchantsClient struct {
	client  *httpdata.HttpClient
	connect sync.Once
}

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

var singletonMerchCli = &MerchantsClient{}

func GetMerchantsClient() *MerchantsClient {
	singletonMerchCli.connect.Do(func() {
		conf := config.Get().Integration.Rest.Merchants
		singletonMerchCli.client = httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized)
	})

	return singletonMerchCli
}
