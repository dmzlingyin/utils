package misc

import "testing"

func TestDistance(t *testing.T) {
	// 北京南站
	lat1 := 39.87
	lng1 := 116.38
	// 清华大学
	lat2 := 40.00
	lng2 := 116.33

	t.Log(Distance(lat1, lng1, lat2, lng2))
}
