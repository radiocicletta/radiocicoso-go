package main


type Programmi struct {
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
}


type Schedule struct {
	Programmi []Programmi `json:"programmi"`
	Adesso struct {
		Start []string `json:"start"`
		End []string `json:"end"`
		ID int `json:"id"`
		Title string `json:"title"`
	} `json:"adesso"`
}


type SortedProgrammi []Programmi


func (s SortedProgrammi) Len() int {
    return len(s) 
}

func (s SortedProgrammi) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s SortedProgrammi) Less(i, j int) bool {
    a := s[i].Start
    b := s[j].Start

    a_hour := a[1].(float64)
    b_hour := b[1].(float64)

    a_minute := a[2].(float64)
    b_minute := b[2].(float64)

    return a_hour < b_hour || a_hour == b_hour && a_minute < b_minute
}
