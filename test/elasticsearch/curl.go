package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"simplest_script/core/svc"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/distanceunit"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/geodistancetype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

type GoodsInfo struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Price    float64  `json:"price"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func UnitTest() {
	ctx := context.Background()
	index := "test"

	// 新增
	// doc := GoodsInfo{
	// 	Id:   3,
	// 	Name: "测试商品3",
	// 	Location: Location{
	// 		Lat: 34.0522,
	// 		Lon: -118.2437,
	// 	},
	// 	Price: 18.88,
	// }
	// id := strconv.Itoa(doc.Id)
	// res, err := svc.NewElasticsearch().Create(ctx, index, id, doc)

	// fmt.Printf("res: %v, err: %v", res, err)

	// 查询
	// 1. 构建参考点坐标
	geoPoint := types.LatLonGeoLocation{
		Lat: 40.768,
		Lon: -73.9865,
	}

	// 2. 初始化 GeoDistanceSort map（key=字段名，value=坐标列表）
	geoDistanceMap := make(map[string][]types.GeoLocation)
	geoDistanceMap["location"] = []types.GeoLocation{geoPoint}

	// 3. 转换枚举类型（转为指针类型，匹配结构体定义）
	ignoreUnmapped := true

	sort := types.SortOptions{
		GeoDistance_: &types.GeoDistanceSort{
			GeoDistanceSort: geoDistanceMap,       // 核心：配置字段名和坐标
			Order:           &sortorder.Asc,       // 排序方向（指针类型）
			Unit:            &distanceunit.Meters, // 距离单位（指针类型）
			DistanceType:    &geodistancetype.Arc, // 距离计算方式（指针类型）
			IgnoreUnmapped:  &ignoreUnmapped,      // 是否忽略未映射字段（指针类型）
			Mode:            nil,                  // 非必需字段置空
			Nested:          nil,                  // 非必需字段置空
		},
	}

	res, err := svc.NewElasticsearch().Search(ctx, index, nil, 1, 10, sort)
	str, _ := json.Marshal(res)
	fmt.Printf("res: %v, err: %v", string(str), err)
}
