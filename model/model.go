package model

import "github.com/golang-module/carbon/v2"

type User struct {
	Id            int64
	Name          string
	Password      string
	FollowCount   int64
	FollowerCount int64
}

type Comment struct {
	CommentId   int64
	UserId      int64
	VideoId     int64
	CommentText string
	CreateTime  string
}

type FavouriteUser struct {
	AuthorId    int64
	FavouriteId int64
}

type FavouriteVideo struct {
	UserId    int64
	VideoId   int64
	Favourite int8
}

type Video struct {
	Id             int64
	AuthorId       int64
	PlayUrl        string
	CoverUrl       string
	FavouriteCount int64
	CommentCount   int64
	CreateTime     carbon.DateTime
}
