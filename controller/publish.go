package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"net/http"
	"os"
	"os/exec"
)

type VideoListResponse struct {
	_type.Response
	VideoList []_type.Video `json:"video_list"`
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	c.JSON(http.StatusOK, VideoListResponse{
		Response: _type.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}

func Action(c *gin.Context) {
	token := c.PostForm("token")

	//获取用户上传的视频名称
	title := c.PostForm("title")

	//视频路径
	filepath1 := "./Data/video/" + token + "/"
	_, err := os.Stat(filepath1)
	if err != nil {
		//创建用户token名称的文件夹
		if os.IsNotExist(err) {
			err := os.Mkdir(filepath1, os.ModePerm)
			if err != nil {
				return
			}
		}
	}

	//r.ParseMultipartForm(32 << 20)
	//获取上传的文件
	file, err := c.FormFile("data")

	if err != nil {
		c.String(http.StatusBadRequest, "A BAD REQUEST")
		return
	}

	//最后的具体的视频路径
	titleSum := filepath1 + title + ".mp4"
	//保存文件到本地中
	err = c.SaveUploadedFile(file, titleSum)
	if err != nil {
		fmt.Print("保存视频失败")
	}

	//设置照片的保存路径
	photoSum := "./Data/photo/" + token + "/"
	_, exist := os.Stat(photoSum)
	if exist != nil {
		//创建用户token名称的文件夹
		if os.IsNotExist(exist) {
			err := os.Mkdir(photoSum, os.ModePerm)
			if err != nil {
				return
			}
		}
	}

	photoSum = photoSum + title + ".bmp"

	//调用ffmpeg截图视频并将截图保存至
	cmd := exec.Command("ffmpeg", "-i", titleSum, "-y", "-f", "image2", "-ss", "00:00:02", "-vframes", "1", photoSum)
	err = cmd.Run()
	if err != nil {
		fmt.Print(err)
		fmt.Print("失败")
	}

	//c.String(http.StatusOK, fmt.Sprintf("%s,upload", file.Filename))
	//获取作者编号
	m := Dbm.GerAllUser()
	authorId := m[token].Id
	lastTime := c.PostForm("latest_time")
	if lastTime == "" {
		lastTime = carbon.Now().ToDateTimeString()
	}

	titleSum = titleSum[1:]
	photoSum = photoSum[1:]
	_ = Dbm.InsertVideo(authorId, titleSum, photoSum, lastTime)

	c.JSON(http.StatusOK, _type.Response{
		StatusCode: 0,
		StatusMsg:  title + " uploaded successfully",
	})
}
