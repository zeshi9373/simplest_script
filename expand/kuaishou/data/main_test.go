package data_test

import (
	"simplest_script/expand/kuaishou/data"
	"testing"
)

func TestKuaishouReportData_AccountReport(t *testing.T) {
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
			name:        "normal",
			accessToken: "517c65dc28ee7b7822942ef8497ec0b2",
			param: &data.AccountReportDataRequest{
				AdvertiserID:        84315186,
				StartDate:           "2025-12-17",
				EndDate:             "2025-12-17",
				Page:                1,
				PageSize:            2,
				TemporalGranularity: "HOURLY",
			},
			list:    l,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := data.NewKuaishouReportData()
			_, gotErr := h.AccountReport(tt.accessToken, tt.param, tt.list)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AccountReport() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AccountReport() succeeded unexpectedly")
			}
		})
	}
}
