package logger_test

import (
	"simplest_script/test"
	"simplest_script/test/logger"
	"testing"
)

func TestInfo(t *testing.T) {

	test.Init()

	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "test info",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Info()
		})
	}
}
