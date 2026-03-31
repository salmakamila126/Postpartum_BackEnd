package timeutil

import "time"

var wib *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	wib = loc
}

func NowWIB() time.Time {
	return time.Now().In(wib)
}

func ToWIB(t time.Time) time.Time {
	return t.In(wib)
}

func StartOfDay(t time.Time) time.Time {
	t = t.In(wib)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, wib)
}

func FormatHour(t time.Time) string {
	return t.In(wib).Format("15:04")
}

const (
	DateOnlyFormat = "2006-01-02"
	TimeOnlyFormat = "15:04"
	RFC3339Format  = time.RFC3339
)

func ParseRFC3339(value string) (time.Time, error) {
	return time.ParseInLocation(time.RFC3339, value, wib)
}

func ParseDate(value string) (time.Time, error) {
	return time.ParseInLocation(DateOnlyFormat, value, wib)
}
