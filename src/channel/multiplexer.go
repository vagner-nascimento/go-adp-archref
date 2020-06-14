package channel

func Multiplex(bytesCh ...<-chan interface{}) <-chan interface{} {
	uniqueCh := make(chan interface{})

	go func() {
		totalChannels := len(bytesCh)
		var closedChannels int

		for _, ch := range bytesCh {
			go forwardAny(ch, uniqueCh, &closedChannels)
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

func forwardAny(from <-chan interface{}, to chan interface{}, closedChannels *int) {
	for err := range from {
		to <- err
	}

	*closedChannels++
}
