package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"text/template"

	cfenv "github.com/cloudfoundry-community/go-cfenvnested"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
)

type Message struct {
	ProxyUrl string
}

type FilterApps struct {
	SelfAppID string
	AppID     string
}

func main() {
	fmt.Println("Loading configuration...")
	var elasticBackendURL string
	var elasticProxyURL string
	var selfAppID = ""
	appEnv, enverr := cfenv.Current()
	if enverr != nil {
		elasticBackendURL = "http://localhost:9200"
		elasticProxyURL = "http://localhost:3000/elasticsearch"
	} else {
		logstash, err := appEnv.Services.WithTag("logstash")
		if err == nil {
			hostname := logstash[0].Credentials["hostname"].(string)
			ports := logstash[0].Credentials["ports"].(map[string]interface{})
			elasticSearchPort := ports["9200/tcp"]
			elasticBackendURL = fmt.Sprintf("http://%s:%s", hostname, elasticSearchPort)
		} else {
			log.Fatal("Unable to find elastic search service")
		}

		elasticProxyURL = fmt.Sprintf("//%s/elasticsearch", appEnv.ApplicationURIs[0])
		selfAppID = appEnv.ApplicationID
	}
	fmt.Printf("Starting kibana to backend elastic search %s...\n", elasticBackendURL)
	m := martini.Classic()
	m.Get("/config.js", func() string {
		var buffer bytes.Buffer
		configTmpl, _ := template.New("config.tmpl").ParseFiles("./config.tmpl")
		configTmpl.Execute(&buffer, Message{ProxyUrl: elasticProxyURL})
		return string(buffer.Bytes())
	})
	m.Get("/app/dashboards/apps-logs.json", func() string {
		var buffer bytes.Buffer
		configTmpl, _ := template.New("apps-logs.tmpl").Delims("[{", "}]").ParseFiles("./apps-logs.tmpl")
		configTmpl.Execute(&buffer, FilterApps{SelfAppID: selfAppID})
		return string(buffer.Bytes())
	})
	m.Get("/app/dashboards/app-logs-:app_guid.json", func(params martini.Params) string {
		var buffer bytes.Buffer
		appGUID := params["app_guid"]
		configTmpl, _ := template.New("app-logs.tmpl").Delims("[{", "}]").ParseFiles("./app-logs.tmpl")
		configTmpl.Execute(&buffer, FilterApps{SelfAppID: selfAppID, AppID: appGUID})
		return string(buffer.Bytes())
	})

	elasticsearchProxy := func(w http.ResponseWriter, r *http.Request) {
		proxyURL := elasticBackendURL + "/" + r.URL.Path[15:]
		fmt.Println("proxy:", proxyURL)
			logstash1, err := appEnv.Services.WithTag("logstash")
			if err == nil {
				var username string
				var password string
				var ok bool
				username, ok := logstash1[0].Credentials["username"]
				if username, ok = logstash1[0].Credentials["username"].(string); !ok {
   				fmt.Println("Unsupported non-string value found for ElasticSearch username")
				}
				if password, ok = logstash1[0].Credentials["password"].(string); !ok {
   				fmt.Println("Unsupported non-string value found for ElasticSearch password")
				}
				if username != "" && password != "" {
					fmt.Println("I am setting header")
					elasticsearchAuth := fmt.Sprintf("%s:%s", username, password)
					encodedelasticsearchAuth :=  base64.StdEncoding.EncodeToString([]byte(elasticsearchAuth))
					BasicAuth := fmt.Sprintf("%s %s", "Basic" , encodedelasticsearchAuth)
					r.Header.Set ("Authorization" , BasicAuth)
				}
			}
		remote, err := url.Parse(proxyURL)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.URL.Path = "/"
		proxy.ServeHTTP(w, r)
	}

	// Proxy requests to Elastic Search
	m.Group("/elasticsearch", func(r martini.Router) {
		r.NotFound(elasticsearchProxy)
	})

	user := os.Getenv("KIBANA_USERNAME")
	pass := os.Getenv("KIBANA_PASSWORD")
	if user != "" && pass != "" {
		m.Use(auth.Basic(user, pass))
	}
	m.Run()
}
