package main

import (
	"github.com/Knorkebrot/fatberris/fatberris-lib"
	"fmt"
	"net/http"
	"net/http/fcgi"
	"time"
	"log"
)

type cacheObj struct {
	updated time.Time
	value string
}

var (
	cache map[string]cacheObj = make(map[string]cacheObj, 4)
)

const (
	PAGE string = `<!DOCTYPE html>
<html>
<head>
	<title>Fat Berri's m3u</title>
	<meta name="viewport" content="width=device-width" />
	<style type="text/css">
		body {
			background:	#efefef;
			margin:		0;
		}
		#content {
			background:	#fff;
			font-family:	sans-serif;
			font-size:	16px;
			max-width:	240px;
			margin:		50px auto 0;
			padding:	20px 30px;
			box-shadow:	0px 0px 7px #555;
		}
		p > a, p > a:visited {
			color:		#555;
		}
		small, small a {
			color:		#aaa;
		}
	</style>
</head>
<body>
	<div id="content">
		<p>What's your mood today?</p>
		<p><a href="?mood=chill">Chill</a>
		   <a href="?mood=up">Up</a>
		   <a href="?mood=down">Down</a>
		   <a href="?mood=mix">Mix</a></p>
		<p><small>Streams provided by <a href="http://fatberris.com/">Fat Berri's</a><br>
			  Converter by bo (<a href="http://kbct.de/">kbct.de</a>)<br>
			  Source code available at <a href="https://github.com/Knorkebrot/fatberris">github.com</a></small></p>
	</div>
</body>
</html>
`
)

func output(w http.ResponseWriter, name, out string) {
	w.Header().Set("Content-Type", "audio/x-mpegurl; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=" + name + ".m3u")
	fmt.Fprint(w, out)
}

func handler(w http.ResponseWriter, r *http.Request) {
	mood := r.FormValue("mood")

	if mood == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, PAGE)
		return;
	}

	if c, ok := cache[mood]; ok && time.Since(c.updated).Hours() < 24 {
		log.Printf("cached: %s\n", mood)
		output(w, mood, cache[mood].value)
		return
	}

	out, err := fatberris.GetM3u([]string{mood})
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("fetched: %s\n", mood)

	cache[mood] = cacheObj{time.Now(), out}

	output(w, mood, out)
}

func main() {
	err := fcgi.Serve(nil, http.HandlerFunc(handler))
	//err := http.ListenAndServe("localhost:1234", http.HandlerFunc(handler))
	if err != nil {
		log.Fatal(err)
	}
}
