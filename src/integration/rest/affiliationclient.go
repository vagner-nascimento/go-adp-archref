package rest

import (
	"github.com/vagner-nascimento/go-adp-bridge/config"
	appentity "github.com/vagner-nascimento/go-adp-bridge/src/app/entity"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/httpdata"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-bridge/src/singleton"
)

type AffiliationsClient struct {
	client *httpdata.HttpClient
}

var singAffCli singleton.SingResource

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

func GetAffiliationsClient() *AffiliationsClient {
	singAffCli.Once.Do(func() {
		conf := config.Get().Integration.Rest.Affiliations

		singAffCli.Resource = &AffiliationsClient{
			client: httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized),
		}
	})

	return singAffCli.Resource.(*AffiliationsClient)
}
