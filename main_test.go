package main

import (
	"testing"
)

// func Test_readTestMessage(t *testing.T) {
// 	t.Run("dts", func(t *testing.T) {
// 		fmt.Println(time.Now().Unix() - int64(86400))
// 	})
// }

func Test_reapInstances(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reapInstances()
		})
	}
}
