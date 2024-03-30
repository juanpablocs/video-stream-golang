package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Thumbnails  []Thumbnail        `bson:"thumbnails" json:"thumbnails"`
	UploaderID  primitive.ObjectID `bson:"uploaderId" json:"uploaderId"`
	Duration    float64            `bson:"duration" json:"duration"`
	VttTrack    string             `bson:"vttTrack" json:"vttTrack"` // URL del archivo VTT
	Status      VideoStatus        `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	VideoUpload VideoUpload        `bson:"videoUpload,omitempty" json:"videoUpload,omitempty"`
}

type VideoStatus string

const (
	StatusIdle       VideoStatus = "idle"
	StatusConverting VideoStatus = "converting"
	StatusFinished   VideoStatus = "finished"
	StatusError      VideoStatus = "error"
)

type VttTrack struct {
	Language string `bson:"language" json:"language"`
	Label    string `bson:"label" json:"label"`
	Src      string `bson:"src" json:"src"`
	Default  bool   `bson:"default" json:"default"`
}

type Thumbnail struct {
	Size     string `bson:"size" json:"size"`         // "small", "large"
	Position string `bson:"position" json:"position"` // "start", "middle", "end"
	Path     string `bson:"path" json:"path"`
	Default  bool   `bson:"default" json:"default"`
}

// VideoUpload representa un archivo de video subido por el usuario.
type VideoUpload struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	VideoID          primitive.ObjectID `bson:"videoId" json:"videoId"`   // Referencia al Video
	FileSize         int64              `bson:"filesize" json:"filesize"` // En bytes
	Url              string             `bson:"url" json:"url"`           // URL del archivo de video original
	VideoCodec       string             `bson:"videoCodec" json:"videoCodec"`
	VideoBitRate     int64              `bson:"videoBitrate" json:"videoBitrate"`
	VideoWidth       int                `bson:"videoWidth" json:"videoWidth"`
	VideoHeight      int                `bson:"videoHeight" json:"videoHeight"`
	VideoAspectRatio string             `bson:"videoAspectRatio" json:"videoAspectRatio"`
	AudioCodec       string             `bson:"audioCodec" json:"audioCodec"`
	AudioChannels    int                `bson:"audioChannels" json:"audioChannels"`
	AudioSampleRate  int                `bson:"audioSampleRate" json:"audioSampleRate"`
}

// VideoUploadTranscoded representa un archivo de video transcodificado.
type VideoUploadTranscoded struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	VideoID          primitive.ObjectID `bson:"videoId" json:"videoId"`             // Referencia al Video
	VideoUploadID    primitive.ObjectID `bson:"videoUploadId" json:"videoUploadId"` // Referencia al VideoUpload original
	FileSize         int64              `bson:"fileSize" json:"fileSize"`           // En bytes
	UrlMP4           string             `bson:"urlMP4" json:"urlMP4"`
	UrlHLS           string             `bson:"urlHLS" json:"urlHLS"`
	VideoCodec       string             `bson:"videoCodec" json:"videoCodec"`
	VideoBitRate     int64              `bson:"videoBitrate" json:"videoBitrate"`
	VideoWidth       int                `bson:"videoWidth" json:"videoWidth"`
	VideoHeight      int                `bson:"videoHeight" json:"videoHeight"`
	VideoAspectRatio string             `bson:"videoAspectRatio" json:"videoAspectRatio"`
	AudioCodec       string             `bson:"audioCodec" json:"audioCodec"`
	AudioChannels    int                `bson:"audioChannels" json:"audioChannels"`
	AudioSampleRate  int                `bson:"audioSampleRate" json:"audioSampleRate"`
}

type VideoRequest struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	UploaderID  primitive.ObjectID `json:"uploaderId"`
	Url         string             `json:"url"`
}
