package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
    "regexp"
    "math/rand"
)

const (
	nickname string = "radiocicoso"
	server   string = "irc.freenode.net:6667"
)

var channels = []string{"#radiocicletta"}
var goodguys = []string{"leonardo", "Cassapanco", "ineff", "autoscatto", "Biappi"}

func handlerHelp(e *irc.Event) (string, string) {
    var replyto string 

	if e.Arguments[0] == nickname {
		replyto = e.Nick
	} else {
		replyto = e.Arguments[0]
	}
	return "Per chiedermi qualcosa, usa un comando" +
           " preceduto da una chiocciola (ad" +
           " es. @help). per conoscere la descrizione di tutti i comandi" +
           " scrivi @longhelp\nComandi disponibili:" +
           " @help @longhelp @cosera @oggi @inonda @ascolto @podcast",
		replyto
}

func handlerLonghelp(e *irc.Event) (string, string) {
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

func handlerAscolto(e *irc.Event) (string, string) {
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

func handlerPodcast(e *irc.Event) (string, string) {
	var replyto string
	if e.Arguments[0] == nickname {
		replyto = e.Nick
	} else {
		replyto = e.Arguments[0]
	}

	return GetMixcloudPodcasts(), replyto
}

func handlerInonda(e *irc.Event) (string, string) {
	var replyto string

	if e.Arguments[0] == nickname {
		replyto = e.Nick
	} else {
		replyto = e.Arguments[0]
	}

	return GetNowOnair(), replyto
}

func handlerOggi(e *irc.Event) (string, string) {
	var replyto string

	if e.Arguments[0] == nickname {
		replyto = e.Nick
	} else {
		replyto = e.Arguments[0]
	}

	return GetTodaySchedule(), replyto
}

func handlerCosera(e *irc.Event) (string, string) {
	var replyto string

	if e.Arguments[0] == nickname {
		replyto = e.Nick
	} else {
		replyto = e.Arguments[0]
	}

	return GetStreamMetadata(), replyto
}

func handlerHello(e *irc.Event, a ...string) (string, string){
	var replyto string

    answers := []string{
        "o/",
        "hello",
        "ciao!",
        "servus",
        "ciaociao",
        "cia'",
    }

	if e.Arguments[0] == nickname {
		replyto = e.Nick
	} else {
		replyto = e.Arguments[0]
	}

    return answers[rand.Intn(len(answers))], replyto
}


func handlerNoembed(e *irc.Event, args ...string) (string, string){
	var replyto string

	if e.Arguments[0] == nickname {
		replyto = e.Nick
	} else {
		replyto = e.Arguments[0]
	}

    if len(args) > 0{ 
        uri := args[0]
        reply := GetNoembedData(uri)
        return reply, replyto
    }

    return "", replyto

}

type IrcReaction struct {
    Pattern *regexp.Regexp
    Handler func(*irc.Event, ...string) (string, string)
}

func (i *IrcReaction) Match(pattern string) bool {
    return i.Pattern.MatchString(pattern)
}

func main() {
	ircconn := irc.IRC(nickname, nickname)
	cmdqueryhandlers := map[string]func(*irc.Event) (string, string){
        "@help"    : handlerHelp,
		"@longhelp": handlerLonghelp,
		"@ascolto":  handlerAscolto,
		"@podcast":  handlerPodcast,
		"@cosera":   handlerCosera,
		"@inonda":   handlerInonda,
		"@oggi":     handlerOggi,
	}

    reactions := []IrcReaction{
        IrcReaction{regexp.MustCompile("https?://[^ ]+"), handlerNoembed},
        IrcReaction{regexp.MustCompile("ciao|buon(giorno|di|asera|anotte)|hello"), handlerHello},
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
			} else {
                for _, rx := range reactions {
                    if rx.Match(k) {
                        msg, replyto = rx.Handler(event, k)
                        for _, m := range strings.Split(msg, "\n") {
                            if m != "" {
                                ircconn.Privmsgf(replyto, m)
                            }
                        }
                        break
                    }
                }
            }
		}
	})

	ircconn.Connect(server)
	for _, channel := range channels {
		ircconn.Join(channel)
	}
	ircconn.Loop()
}
