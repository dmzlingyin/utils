package backoff

func Retry(count int, fn func() error) error {
	for i := 0; i < count; i++ {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
