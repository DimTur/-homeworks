package main

type APIResponse struct {
	Id          string `json:"id"`
	Author      string `json:"author"`
	Width       int64  `json:"width"`
	Height      int64  `json:"height"`
	Url         string `json:"url"`
	DownloadUrl string `json:"download_url"`
}

type DownloadTask struct {
	Id          string
	DownloadUrl string
}
