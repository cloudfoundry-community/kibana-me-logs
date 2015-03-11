package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"text/template"

	cfenv "github.com/cloudfoundry-community/go-cfenvnested"
	"github.com/go-martini/martini"
)

type message struct {
	url string
}

func main() {
	var elasticURL string
	appEnv, enverr := cfenv.Current()
	if enverr != nil {
		elasticURL = "http://localhost:9200"
	} else {
		elasticSearch, err := appEnv.Services.WithTag("elasticsearch")
		fmt.Println(appEnv)
		if err == nil {
			u, _ := url.Parse(elasticSearch[0].Credentials["uri"].(string))
			password, _ := u.User.Password()
			elasticURL = fmt.Sprintf("http://%s/api-key/%s", u.Host, password)
		} else {
			log.Fatal("Unable to find elastic search service")
		}
	}
	m := martini.Classic()
	m.Get("/config.js", func() string {
		var buffer bytes.Buffer
		configTmpl, _ := template.New("config.tmpl").ParseFiles("./config.tmpl")
		configTmpl.Execute(&buffer, message{url: elasticURL})
		return string(buffer.Bytes())
	})
	m.Run()
}
