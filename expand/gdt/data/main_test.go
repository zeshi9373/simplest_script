package data_test

import (
	"simplest_script/expand/gdt/data"
	"testing"
)

func TestGdtReportData_AccountReport(t *testing.T) {
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
			accessToken: "590ec9e4c1cb6ae36f30d2913b9d40ee",
			param: &data.AccountReportDataRequest{
				AccountID: "72020110",
				DateRange: data.DateRange{
					StartDate: "2025-12-10",
					EndDate:   "2025-12-10",
				},
				Fields:   []string{"cost", "hour"},
				GroupBy:  []string{"hour"},
				Level:    "REPORT_LEVEL_ADVERTISER",
				Page:     1,
				PageSize: 10,
			},
			list:    l,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := data.NewGdtReportData()
			gotErr := h.AccountReport(tt.accessToken, tt.param, tt.list)
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
