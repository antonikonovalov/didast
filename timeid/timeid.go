package timeid

import "time"

func New() int64 {
	return time.Now().UTC().UnixNano()
}

func DetectInterval(id int64, interval time.Duration) []int64 {
	t := time.Unix(0, id).Round(interval)
	return []int64{t.Add(-interval).UTC().UnixNano(), t.UTC().UnixNano()}
}
