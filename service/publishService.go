package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/Dao"
	_type "github.com/RaymondCode/simple-demo/type"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"net/http"
	"os"
	"os/exec"
)

func PublishListService(userId string) []_type.Video {
	videolist := Dao.Vdi.GetUserPublish(userId)
	return videolist
}

func ActionService(c *gin.Context, token string, title string) bool {
	//视频路径
	filepath1 := "./Data/video/" + token + "/"
	_, err := os.Stat(filepath1)
	if err != nil {
		//创建用户token名称的文件夹
		if os.IsNotExist(err) {
			err := os.Mkdir(filepath1, os.ModePerm)
			if err != nil {
				return false
			}
		}
	}
	//获取上传的文件
	file, err := c.FormFile("data")

	if err != nil {
		c.String(http.StatusBadRequest, "A BAD REQUEST")
		return false
	}

	//最后的具体的视频路径
	titleSum := filepath1 + title + ".mp4"
	//保存文件到本地中
	err = c.SaveUploadedFile(file, titleSum)
	if err != nil {
		fmt.Print("保存视频失败")
		return false
	}

	//设置照片的保存路径
	photoSum := "./Data/photo/" + token + "/"
	_, exist := os.Stat(photoSum)
	if exist != nil {
		//创建用户token名称的文件夹
		if os.IsNotExist(exist) {
			err := os.Mkdir(photoSum, os.ModePerm)
			if err != nil {
				return false
			}
		}
	}
	photoSum = photoSum + title + ".bmp"
	//调用ffmpeg截图视频并将截图保存至
	cmd := exec.Command("ffmpeg", "-i", titleSum, "-y", "-f", "image2", "-ss", "00:00:02", "-vframes", "1", photoSum)
	err = cmd.Run()
	if err != nil {
		fmt.Println("截图时失败，err =", err)
		return false
	}
	//c.String(http.StatusOK, fmt.Sprintf("%s,upload", file.Filename))
	//获取作者编号
	m := Dao.Udi.GerAllUser()
	authorId := m[token].Id
	lastTime := c.PostForm("latest_time")
	if lastTime == "" {
		lastTime = carbon.Now().ToDateTimeString()
	}
	titleSum = titleSum[1:]
	photoSum = photoSum[1:]
	_ = Dao.Vdi.InsertVideo(authorId, titleSum, photoSum, lastTime)
	return true
}
