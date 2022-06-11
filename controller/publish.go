package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
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
	m := dbm.GerAllUser()
	authorId := m[token].Id
	lastTime := c.PostForm("latest_time")
	if lastTime == "" {
		lastTime = carbon.Now().ToDateTimeString()
	}

	titleSum = titleSum[1:]
	photoSum = photoSum[1:]
	_ = dbm.InsertVideo(authorId, titleSum, photoSum, lastTime)

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  title + " uploaded successfully",
	})
}
