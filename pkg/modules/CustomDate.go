package modules

import (
	"strings"
	"time"
)

type CustomDate time.Time

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s) // Это эталонный формат даты в Go
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}
