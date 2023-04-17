package adapters

import (
	"go_another_chatgpt/repo/censor"
	"reflect"
	"testing"
)

func TestLocalDirtyFilter_MakeTextAuditing(t *testing.T) {
	type args struct {
		id   string
		text string
	}
	tests := []struct {
		name       string
		args       args
		wantResult *censor.TextAuditingResult
		wantErr    bool
	}{
		// TODO: Add test cases.
		{"test01", args{"001", "我是习近平"}, &censor.TextAuditingResult{false, "我是***"}, false},
		{"test01", args{"002", "我是*习近平"}, &censor.TextAuditingResult{false, "我是****"}, false},
		{"test01", args{"002", "大肉棒"}, &censor.TextAuditingResult{false, "大**"}, false},
		{"test01", args{"002", "共产党"}, &censor.TextAuditingResult{true, "共产党"}, false},
		{"test01", args{"002", "今天天气不错"}, &censor.TextAuditingResult{true, "今天天气不错"}, false},
		{"test01", args{"002", "https://pornhub.com"}, &censor.TextAuditingResult{false, "https://***********"}, false},
		{"test01", args{"002", "mypornhub.com"}, &censor.TextAuditingResult{false, "my***********"}, false},
	}
	r := NewLocalDirtyFilter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := r.MakeTextAuditing(tt.args.id, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeTextAuditing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("MakeTextAuditing() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
