package elasticsearch_test

import (
	"simplest_script/test"
	"simplest_script/test/elasticsearch"
	"testing"
)

func TestUnitTest(t *testing.T) {
	test.Init()
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "UnitTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elasticsearch.UnitTest()
		})
	}
}
