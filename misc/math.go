package misc

import (
	"math"
)

// EarthRadius 地球半径，单位米
const EarthRadius = 6367000

// Distance 计算两个经纬度之间的距离（单位：米）
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// 将纬度和经度从度转换为弧度
	lat1Rad := toRadians(lat1)
	lon1Rad := toRadians(lon1)
	lat2Rad := toRadians(lat2)
	lon2Rad := toRadians(lon2)

	// 计算纬度和经度的差值
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// 计算Haversine公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 计算并返回距离
	return EarthRadius * c
}

// toRadians 将角度转换为弧度
func toRadians(degree float64) float64 {
	return degree * math.Pi / 180
}
