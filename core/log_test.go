package core

import (
	"testing"

	"github.com/goextension/log"
)

func TestInitGlobalLogger(t *testing.T) {
	type args struct {
		debug bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				debug: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitGlobalLogger(true)
			log.Debugw("test output")
		})
	}
}
