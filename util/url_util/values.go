package url_util

import (
	"net/url"
	"reflect"
	"strings"
)

func ValuesFromObj(itf interface{}) (values *url.Values) {
	values = &url.Values{}
	t := reflect.TypeOf(itf)
	v := reflect.ValueOf(itf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		tags := strings.Split(tag, ",")
		if tags[0] == "" {
			continue
		}
		if len(tags) > 1 && tags[1] == "omitempty" && v.Field(i).IsZero() {
			continue
		}
		values.Set(tags[0], v.Field(i).String())
	}
	return
}
