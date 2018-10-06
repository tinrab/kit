package retry

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	err := Retry(4, func(attempt int) (err error) {
		if attempt < 3 {
			err = errors.New("")
		}
		return err
	})
	if err != nil {
		t.Fail()
	}
}

func TestRetryError(t *testing.T) {
	err := Retry(4, func(attempt int) error {
		return errors.New("")
	})
	if err == nil {
		t.Fail()
	}
}

func TestRetrySleep(t *testing.T) {
	startTime := time.Now().UnixNano()
	sleepTime := 200 * time.Millisecond

	err := Sleep(2, sleepTime, func(attempt int) (err error) {
		if attempt == 0 {
			err = errors.New("")
		}
		return err
	})

	dt := (time.Now().UnixNano() - startTime - int64(sleepTime)) / (int64(time.Millisecond) / int64(time.Nanosecond))
	maxErr := int64(50)
	if err != nil || (dt < -maxErr || dt > maxErr) {
		t.Fail()
	}
}
