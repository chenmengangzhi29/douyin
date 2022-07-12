package service

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
	"github.com/chenmengangzhi29/douyin/pkg/logger"
	"github.com/chenmengangzhi29/douyin/pkg/oss"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type PublishService struct {
	ctx context.Context
}

// NewPublishService new PublishService
func NewPublishService(ctx context.Context) *PublishService {
	return &PublishService{ctx: ctx}
}

// Publish upload video info
func (s *PublishService) Publish(req *publish.PublishActionRequest) error {
	video := req.Data
	title := req.Title
	currentId := req.UserId

	id := time.Now().Unix()
	fileName := strconv.Itoa(int(id)) + ".mp4"

	//将视频保存到本地文件夹
	filePath := oss.Path + "/public/" + fileName
	err := oss.PublishVideoToPublic(video, filePath)
	if err != nil {
		return err
	}

	//分片将视频流上传到oss
	objectKey := "video/" + fileName
	err = oss.PublishVideoToOss(objectKey, filePath)
	if err != nil {
		return err
	}

	//获取视频播放地址
	signedURL, err := oss.QueryOssVideoURL(objectKey)
	if err != nil {
		return err
	}
	videoUrl := strings.Split(signedURL, "?")[0]

	//获取封面
	coverName := strconv.Itoa(int(id)) + ".jpg"
	coverData, err := getSnapshot(videoUrl)
	if err != nil {
		return err
	}

	//上传封面到oss
	objectKey = "cover/" + coverName
	coverReader := bytes.NewReader(coverData)
	err = oss.PublishCoverToOss(objectKey, coverReader)
	if err != nil {
		return err
	}

	//获取封面链接
	signedURL, err = oss.QueryOssCoverURL(objectKey)
	if err != nil {
		return err
	}
	coverUrl := strings.Split(signedURL, "?")[0]

	//将视频信息上传到mysql
	videoRaw := &db.VideoRaw{
		UserId:        currentId,
		Title:         title,
		PlayUrl:       videoUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		UpdatedAt:     time.Now(),
	}

	err = db.PublishVideoData(s.ctx, videoRaw)
	if err != nil {
		return err
	}
	return nil
}

func getSnapshot(videoUrl string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoUrl).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buffer, os.Stdout).
		Run()
	if err != nil {
		logger.Fatal("生成缩略图失败 ", err)
		return nil, err
	}

	img, err := imaging.Decode(buffer)
	if err != nil {
		logger.Fatal("生成缩略图失败 ", err)
		return nil, err
	}

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)

	return buf.Bytes(), err
}
