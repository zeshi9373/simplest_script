package data_test

import (
	"simplest_script/expand/ocean/data"
	"testing"
)

func TestOceanReportData_Report(t *testing.T) {
	l := &data.Data{}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		accessToken string
		param       *data.AccountReportDataRequest
		list        *data.Data
		wantErr     bool
	}{
		{
			name:        "test",
			accessToken: "ef8afbd1677e3baa696629d2068792eba1d1d0ff",
			param: &data.AccountReportDataRequest{
				AdvertiserID: "1850549219738952",
				Dimensions:   []string{"stat_time_hour"},
				StartTime:    "2025-12-10 00:00:00",
				EndTime:      "2025-12-10 16:00:00",
				Filters:      []interface{}{},
				Metrics:      []string{"stat_cost", "attribution_convert_cnt", "show_cnt", "click_cnt"},
				OrderBy:      []interface{}{},
				Page:         "1",
				PageSize:     "100",
			},
			list:    l,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := data.NewOceanReportData()
			gotErr := h.AccountReport(tt.accessToken, tt.param, tt.list)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Report() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Report() succeeded unexpectedly")
			}
		})
	}
}
