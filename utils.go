package dot

import "time"

func Retry(fn func() error, count int, sleep time.Duration) (err error) {
	for i := 0; i < count; i++ {
		if err = fn(); err == nil {
			break
		}

		time.Sleep(sleep)
	}

	return
}
