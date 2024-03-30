package utils

import (
	"reflect"
	"testing"

	"github.com/juanpablocs/video-stream-golang/internal/models"
)

func TestFilterResolutions(t *testing.T) {
	type args struct {
		videoWidth  int
		videoHeight int
	}

	tests := []struct {
		name string
		args args
		want []models.Resolution
	}{
		{
			name: "should return all resolutions for a 1920x1080 video",
			args: args{
				videoWidth:  1920,
				videoHeight: 1080,
			},
			want: []models.Resolution{
				{Width: 1920, Height: 1080, Label: "1080p"},
				{Width: 1280, Height: 720, Label: "720p"},
				{Width: 854, Height: 480, Label: "480p"},
				{Width: 640, Height: 360, Label: "360p"},
				{Width: 426, Height: 240, Label: "240p"},
			},
		},
		{
			name: "should return resolutions for a 320x240 video",
			args: args{
				videoWidth:  320,
				videoHeight: 240,
			},
			want: []models.Resolution{
				{Width: 426, Height: 240, Label: "240p"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterResolutions(tt.args.videoWidth, tt.args.videoHeight); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterResolutions() = %v, want %v", got, tt.want)
			}
		})
	}
}
