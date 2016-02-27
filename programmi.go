package main

type Schedule struct {
	Programmi []struct {
		Start []interface{} `json:"start"`
		BlogID int `json:"blog_id"`
		End []interface{} `json:"end"`
		BlogURL string `json:"blog_url"`
		Title string `json:"title"`
		Stato string `json:"stato"`
		Logo struct {
			URL string `json:"url"`
			Descr string `json:"descr"`
			Title string `json:"title"`
		} `json:"logo"`
		ID int `json:"id"`
	} `json:"programmi"`
	Adesso struct {
		Start []string `json:"start"`
		End []string `json:"end"`
		ID int `json:"id"`
		Title string `json:"title"`
	} `json:"adesso"`
}
