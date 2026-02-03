package workflow

import (
	"fmt"
	"strconv"
	"time"
)

func secondsSince(epoch int64) int64 {
	now := time.Now().Unix()
	return now - epoch
}

func parseInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid int %q", s)
	}
	return v, nil
}

func parseInt64(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int64 %q", s)
	}
	return v, nil
}
