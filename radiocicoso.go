package main

import (
    "github.com/thoj/go-ircevent"
    "strings"
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

//func muori(e *irc.Event) (string, string){
//}



func main() {
    ircconn := irc.IRC(nickname, nickname)
    cmdqueryhandlers := map[string]func(*irc.Event) (string, string){
        "@longhelp": handlerLonghelp,
        "@ascolto": handlerAscolto,
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
