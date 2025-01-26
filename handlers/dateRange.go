package handlers

import "time"

var (
	start time.Time
	end   time.Time
)

func GetDateRange(dateRange string) (start, end time.Time) {
	now := time.Now()

	switch dateRange {
	case "this_month":
		start = time.Date(
			now.Year(),
			now.Month(),
			1,
			0,
			0,
			0,
			0,
			now.Location(),
		)
		end = start.AddDate(0, 1, -1)
	case "last_month":
		start = time.Date(
			now.Year(),
			now.Month()-1,
			1,
			0,
			0,
			0,
			0,
			now.Location(),
		)
		end = start.AddDate(0, 1, -1)
	case "last_6_months":
		start = now.AddDate(0, -6, 0)
		end = now
	case "this_year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	case "last_year":
		start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year()-1, 12, 31, 23, 59, 59, 0, now.Location())
	default:
		start = time.Date(
			now.Year(),
			now.Month(),
			1,
			0,
			0,
			0,
			0,
			now.Location(),
		)
		end = start.AddDate(0, 1, -1)
	}

	return start, end
}
