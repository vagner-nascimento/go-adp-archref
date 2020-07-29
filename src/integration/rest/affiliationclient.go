package restintegration

import (
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/data/httpdata"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"sync"
)

type AffiliationsClient struct {
	client  *httpdata.HttpClient
	connect sync.Once
}

func (mc *AffiliationsClient) GetAffiliation(id string) (affiliation appentity.Affiliation, err error) {
	var (
		status int
		data   []byte
	)

	status, data, err = mc.client.Get("", id)

	if err = handleResponse(status, err, data, "affiliation"); err != nil {
		logger.Error("error on try to get affiliation", err)
	} else {
		affiliation, err = appentity.NewAffiliation(data)
	}

	return
}

var singletonAffCli = &AffiliationsClient{}

func GetAffiliationsClient() *AffiliationsClient {
	singletonAffCli.connect.Do(func() {
		conf := config.Get().Integration.Rest.Affiliations
		singletonAffCli.client = httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized)
	})

	return singletonAffCli
}
