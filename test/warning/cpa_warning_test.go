package warning_test

import (
	"simplest_script/test"
	"simplest_script/test/warning"
	"testing"
)

func TestWarning(t *testing.T) {
	test.Init()
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "TestWarning",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning.Warning()
		})
	}
}

func TestSendWarning(t *testing.T) {
	test.Init()
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "TestSendWarning",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning.SendWarning()
		})
	}
}
