package rest

import (
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/httpdata"
)

type AffiliationsClient struct {
	client *httpdata.HttpClient
}

var singAffCli singletonResource

func (mc *AffiliationsClient) GetAffiliation(id string) (int, []byte, error) {
	return mc.client.Get("", id)
}

func GetAffiliationsClient() *AffiliationsClient {
	singAffCli.once.Do(func() {
		conf := config.Get().Integration.Rest.Affiliations

		singAffCli.resource = &AffiliationsClient{
			client: httpdata.NewHttpClient(conf.BaseUrl, conf.TimeOut, conf.RejectUnauthorized),
		}
	})

	return singAffCli.resource.(*AffiliationsClient)
}
