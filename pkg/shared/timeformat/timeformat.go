package timeformat

import (
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/en"
	"github.com/olebedev/when/rules/common"
)

func ParseTime(timeStr string) (time.Time, error) {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	parsed, err := w.Parse(timeStr, time.Now())
	if err != nil {
		return time.Time{}, err
	}

	return parsed.Time, nil
}