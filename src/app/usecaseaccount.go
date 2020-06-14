package app

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
					mAcc []merchantAccount
					aff  *affiliation
				)

				for ent := range enrichCh {
					switch ent.(type) {
					case []merchantAccount:
						mAcc = ent.([]merchantAccount)
					case affiliation:
						aff = &affiliation{}
						*aff = ent.(affiliation)
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
					mAcc *merchantAccount
					mer  *merchant
				)

				for ent := range enrichCh {
					switch ent.(type) {
					case merchantAccount:
						mAcc = &merchantAccount{}
						*mAcc = ent.(merchantAccount)
					case merchant:
						mer = &merchant{}
						*mer = ent.(merchant)
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
