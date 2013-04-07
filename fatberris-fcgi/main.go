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
	<style type="text/css">
		body {font-family: sans-serif;}
	</style>
</head>
<body>
	<p>What's your mood today?</p>
	<p><a href="?mood=chill">Chill</a>
	   <a href="?mood=up">Up</a>
	   <a href="?mood=down">Down</a>
	   <a href="?mood=mix">Mix</a></p>
	<p><small>Streams provided by <a href="http://fatberris.com/">Fat Berri's</a><br>
	          Converter by bo (<a href="http://kbct.de/">kbct.de</a>)<br>
		  Sourcecode available at <a href="https://github.com/Knorkebrot/fatberris">github.com</a></small></p>
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

	out, err := fatberris.GetMoods([]string{mood})
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
