package main

import (
    "github.com/thoj/go-ircevent"
    "strings"
    "net/http"
    "encoding/json"
    "io"
    "fmt"
    "strconv"
    "regexp"
    "time"
)


const(
    nickname string = "radiocicoso"
    server string = "irc.freenode.net:6667"
)

var channels = []string{"#radiocicletta"}
var goodguys = []string{"leonardo", "Cassapanco", "ineff", "autoscatto", "Biappi"}


func handlerLonghelp(e *irc.Event) (string, string){
    return "Lista dei comandi disponibili:\n \n" +
           "@help      : mostra il messaggio di aiuto\n" +
           "@longhelp  : mostra tutti i comandi disponibili\n" +
           "@cosera    : mostra le informazioni sul brano appena passato\n" +
           "@oggi      : mostra i programmi in onda oggi in radio\n" +
           "@inonda    : mostra il programma ora in onda\n" +
           "@ascolto   : come fare per ascoltare radiocicletta\n" +
           "@podcast   : elenca gli ultimi 5 podcast\n",
           e.Nick
}

func handlerAscolto(e *irc.Event) (string, string){
    var replyto string

    if e.Arguments[0] == nickname {
        replyto = e.Nick
    } else {
        replyto = e.Arguments[0]
    }

    return "Puoi ascoltare radiocicletta in diversi modi:\n" +
           "• Dal tuo browser, collegandoti al sito " +
           "http://www.radiocicletta.it e usando il player del sito\n" +
           "• Usando il tuo programma preferito (VLC, iTunes, RealPlayer...) " +
           "inserendo nella playlist l'indirizzo " +
           "http://stream.radiocicletta.it/stream\n",
           replyto
}


func handlerPodcast(e *irc.Event) (string, string){
    var replyto string
    var message = make([]string, 5, 5)
    var jsondata MixcloudPodcast

    if e.Arguments[0] == nickname {
        replyto = e.Nick
    } else {
        replyto = e.Arguments[0]
    }

    resp, err := http.Get("http://api.mixcloud.com/radiocicletta/cloudcasts/?limit=5")

    defer resp.Body.Close()

    if err != nil {
        return "In questo momento non posso elencare i podcast. Prova più tardi :)",
                replyto
    }

    decoder := json.NewDecoder(resp.Body)
    decoder.Decode(&jsondata)
    for idx, i := range jsondata.Data {
        message[idx] = " • " + i.Name + "\n" + i.Url + "\n"
    }
    return strings.Join(message, "\n"), replyto

}


func handlerInonda(e *irc.Event) (string, string){
    var replyto string
    var jsondata Schedule
    var days = []string{"do", "lu", "ma", "me", "gi", "ve", "sa"}
    var now = time.Now()

    if e.Arguments[0] == nickname {
        replyto = e.Nick
    } else {
        replyto = e.Arguments[0]
    }

    resp, err := http.Get("http://www.radiocicletta.it/programmi.json")

    defer resp.Body.Close()

    if err != nil {
        return "Ora su due piedi non saprei :(", 
                replyto
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

        if startday == dow &&
            (starthour < hour || 
                (starthour == hour && (minute > startminute))) &&
            (starthour > hour || 
                (starthour == hour && (minute < startminute))) {
                    return fmt.Sprintf("Ora in onda: %s", i.Title), replyto
        }
    }
    return "Non saprei :(", replyto

}



func handlerOggi(e *irc.Event) (string, string){
    var replyto string
    var jsondata Schedule
    var days = []string{"do", "lu", "ma", "me", "gi", "ve", "sa"}
    var now = time.Now()

    if e.Arguments[0] == nickname {
        replyto = e.Nick
    } else {
        replyto = e.Arguments[0]
    }

    resp, err := http.Get("http://www.radiocicletta.it/programmi.json")

    defer resp.Body.Close()

    if err != nil {
        return "Ora su due piedi non saprei :(", 
                replyto
    }

    decoder := json.NewDecoder(resp.Body)
    decoder.Decode(&jsondata)

    dow := days[now.Weekday()]

    today := make([]string, 24) // like, 24 hours a day

    j := 0
    for _, i := range jsondata.Programmi {
        startday := i.Start[0].(string)

        if startday == dow {
            today[j] = fmt.Sprintf("%02d:%02d %s",
                int(i.Start[1].(float64)),
                int(i.Start[1].(float64)), 
                i.Title,
            )
            j = j + 1
        }
    }
    return strings.Join(today, "\n"), replyto

}


func handlerCosera(e *irc.Event) (string, string){
    var reply, replyto string

    if e.Arguments[0] == nickname {
        replyto = e.Nick
    } else {
        replyto = e.Arguments[0]
    }

    var client *http.Client = &http.Client{}
    var req, _ = http.NewRequest("GET", "http://stream.radiocicletta.it/stream", nil)
    req.Header.Add("Icy-MetaData", "1")
    resp, _ := client.Do(req)

    if header := resp.Header.Get("Icy-Metaint"); header != "" {

        databytes, _ := strconv.ParseInt(header, 10, 32)
        paddedbytes := databytes + 255

        data := io.LimitReader(resp.Body, paddedbytes) // ensure will read at most #paddedbytes bytes
        buf := make([]byte, paddedbytes, paddedbytes)

        for read, err := 0, error(nil) ; err == nil; read, err = data.Read(buf[read:]) { // metadata
        }

        metadata := fmt.Sprintf("%s", buf[databytes:])

        re := regexp.MustCompile("StreamTitle='([^']*)'")
        titleartist := re.FindStringSubmatch(metadata)
    
        if len(titleartist) >= 2 {
            reply = titleartist[1]
        } else {
            reply = "Non saprei. :("
        }
    } else { 
        reply = "Non saprei. :("
    }

    return reply, 
           replyto
}


func main() {
    ircconn := irc.IRC(nickname, nickname)
    cmdqueryhandlers := map[string]func(*irc.Event) (string, string){
        "@longhelp": handlerLonghelp,
        "@ascolto": handlerAscolto,
        "@podcast": handlerPodcast,
        "@cosera": handlerCosera,
        "@inonda": handlerInonda,
        "@oggi": handlerOggi,
    }

    ircconn.AddCallback("PRIVMSG", func(event *irc.Event) {
        var tokens []string = strings.Fields(event.Message())
        var msg, replyto string

        for _, k := range tokens {
            if handler, ok := cmdqueryhandlers[k]; ok == true {
                msg, replyto = handler(event)
                for _, m := range strings.Split(msg, "\n") { 
                    if m != "" {
                        ircconn.Privmsgf(replyto, m)
                    }
                }
                break
            }
        }
    })

    ircconn.Connect(server)
    for _, channel := range channels {
        ircconn.Join(channel)
    }
    ircconn.Loop()
}
