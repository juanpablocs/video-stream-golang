package usecases

import (
	"context"

	"github.com/juanpablocs/video-stream-golang/internal/models"
	pkgVideo "github.com/juanpablocs/video-stream-golang/pkg/video"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u Usecase) InsertVideoTranscoded(videoID, urlMP4, urlHLS string, videoData *pkgVideo.Metadata) error {
	collection := u.db.Collection("videoUploadsTranscoded")
	videoIDObj, _ := primitive.ObjectIDFromHex(videoID)

	// TODO: migrate to usecase method
	uploadCollection := u.db.Collection("videoUploads")
	var videoUpload *models.VideoUpload
	err := uploadCollection.FindOne(context.TODO(), bson.M{"videoId": videoIDObj}).Decode(&videoUpload)
	if err != nil {
		return err
	}

	newVideo := models.VideoUploadTranscoded{
		ID:               primitive.NewObjectID(),
		VideoID:          videoIDObj,
		VideoUploadID:    videoUpload.ID,
		UrlMP4:           urlMP4,
		UrlHLS:           urlHLS,
		FileSize:         videoData.Video.FileSize,
		VideoCodec:       videoData.Video.Codec,
		VideoWidth:       videoData.Video.Width,
		VideoHeight:      videoData.Video.Height,
		VideoBitRate:     videoData.Video.BitRate,
		VideoAspectRatio: videoData.Video.AspectRatio,
		AudioCodec:       videoData.Audio.Codec,
		AudioChannels:    videoData.Audio.Channels,
		AudioSampleRate:  videoData.Audio.SampleRate,
	}
	_, errInsert := collection.InsertOne(context.TODO(), newVideo)
	if errInsert != nil {
		return errInsert
	}

	return nil
}
