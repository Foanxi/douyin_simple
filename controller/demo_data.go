package controller

import "github.com/RaymondCode/simple-demo/type"

var DemoVideos = []_type.Video{
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "http://10.34.151.198:8080/static/bear.mp4",
		CoverUrl:      "http://10.34.151.198:8080/static/bear-1283347_1280.png",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    true,
	},
}

var DemoComments = []_type.Comment{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = _type.User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
