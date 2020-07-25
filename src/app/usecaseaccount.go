package app

// TODO: receive an interface{} and validate if merchant or seller type
func createAccount(data []byte) (*Account, error) {
	return newAccount(data)
}

func enrichAccount(acc *Account, repo AccountDataHandler) <-chan error {
	resCh := make(chan error)

	go func() {
		defer close(resCh)

		switch acc.Type {
		case getAccountType().merchant:
			{
				enrichCh := getMerchantEnrichmentData(*acc, repo)

				var (
					mAcc []MerchantAccount
					aff  *Affiliation
				)

				for ent := range enrichCh {
					switch ent.(type) {
					case []MerchantAccount:
						mAcc = ent.([]MerchantAccount)
					case Affiliation:
						aff = &Affiliation{}
						*aff = ent.(Affiliation)
					case error:
						resCh <- ent.(error)
					}
				}

				enrichMerchantAccount(acc, mAcc, aff)
			}
		case getAccountType().seller:
			{
				enrichCh := getSellerEnrichmentData(*acc, repo)

				var (
					mAcc *MerchantAccount
					mer  *Merchant
				)

				for ent := range enrichCh {
					switch ent.(type) {
					case MerchantAccount:
						mAcc = &MerchantAccount{}
						*mAcc = ent.(MerchantAccount)
					case Merchant:
						mer = &Merchant{}
						*mer = ent.(Merchant)
					case error:
						resCh <- ent.(error)
					}
				}

				enrichSellerAccount(acc, mer, mAcc)
			}
		}
	}()

	return resCh
}
