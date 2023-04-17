package gin_util

import "reflect"

func InitPagination(p interface{}) {
	eP := reflect.ValueOf(p).Elem()

	ePage := eP.FieldByName("Page")
	eLimit := eP.FieldByName("Limit")

	if ePage.Int() <= 0 {
		ePage.SetInt(1)
	}
	if eLimit.Int() <= 0 {
		eLimit.SetInt(100)
	}
}

type Pagination struct {
	Page      int64
	Limit     int64
	Total     int64
	TotalPage int64

	Data interface{}
}
