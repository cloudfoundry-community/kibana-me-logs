package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	cfenv "github.com/cloudfoundry-community/go-cfenvnested"
	"github.com/go-martini/martini"
)

type Message struct {
	Url string
}

func main() {
	fmt.Println("Loading configuration...")
	var elasticURL string
	appEnv, enverr := cfenv.Current()
	if enverr != nil {
		elasticURL = "http://localhost:9200"
	} else {
		logstash, err := appEnv.Services.WithTag("logstash")
		if err == nil {
			hostname := logstash[0].Credentials["hostname"].(string)
			ports := logstash[0].Credentials["ports"].(map[string]interface{})
			elasticSearchPort := ports["9200/tcp"]
			elasticURL = fmt.Sprintf("http://%s:%s", hostname, elasticSearchPort)
		} else {
			log.Fatal("Unable to find elastic search service")
		}
	}
	fmt.Printf("Starting kibana to backend elastic search %s...\n", elasticURL)
	m := martini.Classic()
	m.Get("/config.js", func() string {
		var buffer bytes.Buffer
		configTmpl, _ := template.New("config.tmpl").ParseFiles("./config.tmpl")
		configTmpl.Execute(&buffer, Message{Url: elasticURL})
		return string(buffer.Bytes())
	})
	m.Run()
}
