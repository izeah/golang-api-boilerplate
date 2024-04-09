package date

import "time"

// Now ...
func Now() *time.Time {
	now := time.Now()
	return &now
}

// NowUTC ...
func NowUTC() *time.Time {
	now := time.Now().UTC()
	return &now
}

// NowLocal ...
func NowLocal() *time.Time {
	now := time.Now().UTC().Add(time.Hour * 7)
	return &now
}

// NowWithLocation ...
func NowWithLocation() *time.Time {
	now := time.Now().In(Location())
	return &now
}

// Location ...
func Location() *time.Location {
	return time.FixedZone("Asia/Jakarta", 7*60*60)
}

func Parse(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, Location())
}

// LastWeek ...
func LastWeek(now time.Time) (start time.Time, end time.Time) {
	end = StartOfWeek(now).Add(-1)

	oneWeek := (24 * 6) * time.Hour
	start = StartOfDay(end.Add(-oneWeek))
	return
}

// LastMonth ...
func LastMonth(now time.Time) (time.Time, time.Time) {
	end := StartOfMonth(now).Add(-time.Nanosecond)
	return StartOfMonth(end), end
}

// StartOfMonth ...
func StartOfMonth(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

// StartOfWeek ...
func StartOfWeek(now time.Time) time.Time {
	wd := now.Weekday()
	if wd == time.Sunday {
		now = now.AddDate(0, 0, -6)
	} else {
		now = now.AddDate(0, 0, -int(wd)+1)
	}
	return StartOfDay(now)
}

// StartOfDay ...
func StartOfDay(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// EndOfDay ...
func EndOfDay(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-1), now.Location())
}
