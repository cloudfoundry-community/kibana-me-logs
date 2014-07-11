package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"text/template"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-martini/martini"
)

type Message struct {
	Url string
}

func main() {
	var elasticUrl string
	appEnv, enverr := cfenv.Current()
	if enverr != nil {
		elasticUrl = "http://localhost:9200"
	} else {
		elasticSearch, err := appEnv.Services.WithTag("elasticsearch")
		fmt.Println(appEnv)
		if err == nil {
			u, _ := url.Parse(elasticSearch[0].Credentials["uri"])
			password, _ := u.User.Password()
			elasticUrl = fmt.Sprintf("http://%s/api-key/%s", u.Host, password)
		} else {
			log.Fatal("Unable to find elastic search service")
		}
	}
	m := martini.Classic()
	m.Get("/config.js", func() string {
		var buffer bytes.Buffer
		configTmpl, _ := template.New("config.tmpl").ParseFiles("./config.tmpl")
		configTmpl.Execute(&buffer, Message{Url: elasticUrl})
		return string(buffer.Bytes())
	})
	m.Run()
}
