package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var videoIdeSequence = GetLastVideoId()

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
	fmt.Println(token)
	//获取用户上传的视频名称
	title := c.PostForm("title")

	filepath1 := "C:/Users/30703/Desktop/video/" + token
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
	//打印日志
	//log.Println(title)
	if err != nil {
		c.String(http.StatusBadRequest, "A BAD REQUEST")
		return
	}
	//保存文件到本地中
	titleSum := filepath1 + "/" + title + ".mp4"
	//保存截图到本地
	//imagePath := filepath1 + "/" + title + ".jpg"
	////获取视频封面图
	//cmd := exec.Command("ffmpeg", "-i", imagePath, "-s", "4cif")
	//buf := new(bytes.Buffer)
	//cmd.Stdout =
	//log.Println(buf)
	//c.SaveUploadedFile(, imagePath)
	c.SaveUploadedFile(file, titleSum)
	c.String(http.StatusOK, fmt.Sprintf("%s,upload", file.Filename))
	//获取作者编号
	m := GerAllUser()
	authorId := m[token].Id
	insertVideo := InsertVideo(authorId, titleSum, "C:/Users/30703/Pictures/桌面壁纸/heihei.jpg")
	log.Println(insertVideo)
}
