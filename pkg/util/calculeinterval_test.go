package util

import "testing"

func TestCalculeInterval(t *testing.T) {
	type args struct {
		duracionMinutos float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "Test 1",
			args: args{duracionMinutos: 1},
			want: 1.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculeInterval(tt.args.duracionMinutos); got != tt.want {
				t.Errorf("CalculeInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
