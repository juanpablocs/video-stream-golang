package video

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type Video struct {
	Codec       string  `json:"codec"`
	BitRate     int64   `json:"bit_rate"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	FileSize    int64   `json:"filesize"`
	Duration    float64 `json:"duration"`
	Filename    string  `json:"filename"`
	AspectRatio string  `json:"aspect_ratio"`
}
type Metadata struct {
	Video Video `json:"video"`
	Audio struct {
		Codec      string `json:"codec"`
		Channels   int    `json:"channels"`
		SampleRate int    `json:"rate"`
	} `json:"audio"`
}

func VideoInfo(videoUrl string, metadata *Metadata) (*Metadata, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", videoUrl)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe error: %w", err)
	}

	// Estructura para analizar la salida de ffprobe
	var ffprobeOutput struct {
		Streams []struct {
			Index       int    `json:"index"`
			CodecName   string `json:"codec_name"`
			Width       int    `json:"width,omitempty"`
			Height      int    `json:"height,omitempty"`
			BitRate     string `json:"bit_rate"`
			AspectRatio string `json:"display_aspect_ratio,omitempty"`
			Duration    string `json:"duration"`
			RFrameRate  string `json:"r_frame_rate,omitempty"`
			SampleRate  string `json:"sample_rate,omitempty"`
			Channels    int    `json:"channels,omitempty"`
			CodecType   string `json:"codec_type"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
			Size     string `json:"size"`
			BitRate  string `json:"bit_rate"`
		} `json:"format"`
	}

	err = json.Unmarshal(out, &ffprobeOutput)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	duration, _ := strconv.ParseFloat(ffprobeOutput.Format.Duration, 64)
	filesize, _ := strconv.ParseInt(ffprobeOutput.Format.Size, 10, 64)
	bitRate, _ := strconv.ParseInt(ffprobeOutput.Format.BitRate, 10, 64)
	metadata.Video.BitRate = bitRate
	metadata.Video.Duration = duration
	metadata.Video.FileSize = filesize

	for _, stream := range ffprobeOutput.Streams {
		switch stream.CodecType {
		case "video":
			if stream.CodecName == "mjpeg" || stream.CodecName == "png" || stream.CodecName == "gif" {
				break
			}
			metadata.Video.Codec = stream.CodecName
			metadata.Video.Width = stream.Width
			metadata.Video.Height = stream.Height
			metadata.Video.AspectRatio = stream.AspectRatio

		case "audio":
			metadata.Audio.Codec = stream.CodecName
			metadata.Audio.Channels = stream.Channels
			sampleRate, _ := strconv.Atoi(stream.SampleRate)
			metadata.Audio.SampleRate = sampleRate
		}

	}
	return metadata, nil
}
