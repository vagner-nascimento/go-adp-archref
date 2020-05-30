package channel

func MultiplexErrors(errsCh ...<-chan error) <-chan error {
	uniqueCh := make(chan error)

	go func() {
		totalChannels := len(errsCh)
		var closedChannels int

		for _, errCh := range errsCh {
			go forwardError(errCh, uniqueCh, &closedChannels)
		}

		for {
			if totalChannels == closedChannels {
				break
			}
		}

		close(uniqueCh)
	}()

	return uniqueCh
}

func forwardError(from <-chan error, to chan error, closedChannels *int) {
	for err := range from {
		to <- err
	}

	*closedChannels++
}
