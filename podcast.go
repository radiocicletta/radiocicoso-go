package main

import (
	"time"
)

type MixcloudPaging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type MixcloudTags struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type MixcloudPictures struct {
	Medium       string `json:"medium"`
	S768wx768h   string `json:"768wx768h"`
	S320wx320h   string `json:"320wx320h"`
	ExtraLarge   string `json:"extra_large"`
	Large        string `json:"large"`
	S640wx640h   string `json:"640wx640h"`
	MediumMobile string `json:"medium_mobile"`
	Small        string `json:"small"`
	S1024wx1024h string `json:"1024wx1024h"`
	Thumbnail    string `json:"thumbnail"`
}

type MixcloudUser struct {
	Url      string           `json:"url"`
	Username string           `json:"username"`
	Name     string           `json:"name"`
	Key      string           `json:"key"`
	Pictures MixcloudPictures `json:"pictures"`
}

type MixcloudData struct {
	Tags          []MixcloudTags   `json:"tags"`
	PlayCount     int              `json:"play_count"`
	User          MixcloudUser     `json:"user"`
	Key           string           `json:"key"`
	CreatedTime   time.Time        `json:"created_time"`
	AudioLength   int              `json:"audio_length"`
	Slug          string           `json:"slug"`
	FavoriteCount int              `json:"favourite_count"`
	ListenerCount int              `json:"listener_count"`
	Name          string           `json:"name"`
	Url           string           `json:"url"`
	Pictures      MixcloudPictures `json:"pictures"`
	RepostCount   int              `json:"repost_count"`
	UpdatedTime   time.Time        `json:"updated_time"`
	CommentCount  int              `json:"comment_count"`
}

type MixcloudPodcast struct {
	Paging MixcloudPaging `json:"paging"`
	Data   []MixcloudData `json:"data"`
	Name   string         `json:"name"`
}
