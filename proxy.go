package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"
)

type commitData struct {
	Before     string
	After      string
	Ref        string
	UserName   string
	Repository struct {
		Url string
	}
}

var (
	listen = flag.String("listen", "localhost:9080", "listen on address")
	logp   = flag.Bool("log", false, "enable logging")
)

func main() {
	flag.Parse()
	proxyHandler := http.HandlerFunc(proxyHandlerFunc)
	log.Fatal(http.ListenAndServe(*listen, proxyHandler))
}

func readerToString(r io.Reader) string {
	if b, err := ioutil.ReadAll(r); err == nil {
		return string(b)
	} else {
		return ""
	}
}

func setGitData(form url.Values, g commitData) {
	form.Set("START", g.Before)
	form.Set("END", g.After)
	form.Set("REFNAME", g.Ref)
	form.Set("URL", g.Repository.Url)
}

func proxyToEndpoint(url string, form url.Values, w http.ResponseWriter) error {
	resp, err := http.PostForm(url, form)
	log.Printf("Posting to: %v\n", url)

	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "ERROR")
	} else {
		defer resp.Body.Close()
		resp.Write(w)
	}
	return err
}

func proxyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if *logp {
		log.Println(r.URL)
	}

	body := readerToString(r.Body)
	decoder := json.NewDecoder(strings.NewReader(body))
	var gitData commitData
	err := decoder.Decode(&gitData)

	log.Printf("Body is: %v\n", body)

	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "JSON body not found!")
	} else if r.FormValue("url") == "" {
		log.Print("URL not found!")
		fmt.Fprintf(w, "URL not found!")
	} else {
		form := make(url.Values)
		setGitData(form, gitData)
		form.Set("payload", body)

		postUrl := r.FormValue("url")
		proxyToEndpoint(postUrl, form, w)
	}
}
