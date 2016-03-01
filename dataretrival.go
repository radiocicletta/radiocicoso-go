package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetMixcloudPodcasts() string {
	var message = make([]string, 5, 5)
	var jsondata MixcloudPodcast

	resp, err := http.Get("http://api.mixcloud.com/radiocicletta/cloudcasts/?limit=5")

	defer resp.Body.Close()

	if err != nil {
		return "In questo momento non posso elencare i podcast. Prova più tardi :)"
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&jsondata)
	for idx, i := range jsondata.Data {
		message[idx] = " • " + i.Name + "\n" + i.Url + "\n"
	}
	return strings.Join(message, "\n")
}

func GetNoembedData(uri string) string {
	var jsondata NoEmbed

    resp, err := http.Get("http://noembed.com/embed?url=" + uri)

	defer resp.Body.Close()

	if err != nil {
		return ""
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&jsondata)

	return  jsondata.Title
}


func GetNowOnair() string {
	var jsondata Schedule
	var days = []string{"do", "lu", "ma", "me", "gi", "ve", "sa"}
	var now = time.Now()

	resp, err := http.Get("http://www.radiocicletta.it/programmi.json")

	defer resp.Body.Close()

	if err != nil {
		return "Ora su due piedi non saprei :("
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&jsondata)

	hour := now.Hour()
	minute := now.Minute()
	dow := days[now.Weekday()]

	for _, i := range jsondata.Programmi {
		startday := i.Start[0].(string)
		starthour := int(i.Start[1].(float64))
		startminute := int(i.Start[2].(float64))

		endhour := int(i.End[1].(float64))
		endminute := int(i.End[2].(float64))

		status := i.Stato

		if startday == dow &&
			status == "1" &&
			(starthour < hour ||
				(starthour == hour && (minute > startminute))) &&
			(endhour > hour ||
				(endhour == hour && (minute < endminute))) {
			return fmt.Sprintf("Ora in onda: %s", i.Title)
		}
	}
	return "Non saprei :("
}

func GetTodaySchedule() string {
	var jsondata Schedule
	var days = []string{"do", "lu", "ma", "me", "gi", "ve", "sa"}
	var now = time.Now()

	resp, err := http.Get("http://www.radiocicletta.it/programmi.json")

	defer resp.Body.Close()

	if err != nil {
		return "Ora su due piedi non saprei :("
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&jsondata)

	dow := days[now.Weekday()]

	todaystr := make([]string, 24) // like, 24 hours a day
	today := make([]Programmi, 24)

	j := 0
	for _, i := range jsondata.Programmi {

		if startday, status := i.Start[0].(string), i.Stato; startday == dow && status == "1" {
			today[j] = i
			j = j + 1
		}
	}
	today = today[:j]
	sort.Sort(SortedProgrammi(today))
	for idx, i := range today {
		todaystr[idx] = fmt.Sprintf("%02d:%02d %s",
			int(i.Start[1].(float64)),
			int(i.Start[2].(float64)),
			i.Title,
		)
	}
	return strings.Join(todaystr, "\n")
}

func GetStreamMetadata() string {
	var client *http.Client = &http.Client{}
	var req, _ = http.NewRequest("GET", "http://stream.radiocicletta.it/stream", nil)
	req.Header.Add("Icy-MetaData", "1")
	resp, _ := client.Do(req)

    defer resp.Body.Close()

	if header := resp.Header.Get("Icy-Metaint"); header != "" {

		databytes, _ := strconv.ParseInt(header, 10, 32)
		paddedbytes := databytes + 255

		data := io.LimitReader(resp.Body, paddedbytes) // ensure will read at most #paddedbytes bytes
		buf := make([]byte, paddedbytes, paddedbytes)

		for read, err := 0, error(nil); err == nil; read, err = data.Read(buf[read:]) { // metadata
		}

		metadata := fmt.Sprintf("%s", buf[databytes:])

		re := regexp.MustCompile("StreamTitle='([^']*)'")
		titleartist := re.FindStringSubmatch(metadata)

		if len(titleartist) >= 2 {
			return titleartist[1]
		} else {
			return "Non saprei. :("
		}
	}
	return "Non saprei. :("
}
