package json_util

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type StrOrUint64 struct {
	Value uint64
}

func (r *StrOrUint64) UnmarshalJSON(data []byte) error {
	sysLpStr := string(data)
	if strings.Contains(sysLpStr, "-nan") {
		r.Value = uint64(math.Inf(-1))
		return nil
	}
	if strings.Contains(sysLpStr, "nan") {
		r.Value = uint64(math.Inf(1))
		return nil
	}
	if strings.Contains(sysLpStr, "\"") {
		sysLpStr = strings.ReplaceAll(sysLpStr, "\"", "")
	}

	return json.Unmarshal([]byte(sysLpStr), &r.Value)
}

type StrOrFloat64 struct {
	Value float64
}

func NanAs(value, def float64) float64 {
	if math.IsNaN(value) {
		return def
	}
	return value
}

func (r *StrOrFloat64) MarshalJSON() (bs []byte, err error) {
	s := ""
	if r.Value == math.Inf(1) {
		s = "nan"
	} else if r.Value == math.Inf(-1) {
		s = "-nan"
	} else {
		s = fmt.Sprintf("\"%v\"", r.Value)
	}
	bs = []byte(s)
	return
}

func (r *StrOrFloat64) UnmarshalJSON(data []byte) error {
	sysLpStr := string(data)
	if strings.Contains(sysLpStr, "-nan") {
		r.Value = math.Inf(-1)
		return nil
	}
	if strings.Contains(sysLpStr, "nan") {
		r.Value = math.Inf(1)
		return nil
	}
	if strings.Contains(sysLpStr, "\"") {
		sysLpStr = strings.ReplaceAll(sysLpStr, "\"", "")
		//if sysLpStr == "-" {
		//	r.Value = math.NaN()
		//	return nil
		//}
		var err error
		r.Value, err = strconv.ParseFloat(sysLpStr, 64)
		if err != nil {
			r.Value = math.NaN()
			return nil
		}
	}

	return json.Unmarshal([]byte(sysLpStr), &r.Value)
}
