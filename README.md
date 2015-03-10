Kibana Me Logs
==============

Draining logs from your Cloud Foundry hosted application to a backend Logstash/Elastic Search is great - it drains them all day long. Except, you can't see them. They are stored in Elastic Search and you have no way to access either Elastic Search nor a Kibana UI to view the logs.

The application hosts Kibana 3 and a proxy that binds to your Logstash/Elastic Search backend service. You can now see your logs in Kibana!

Usage
-----

```
go get github.com/cloudfoundry-community/kibana-me-logs
cd $GOPATH/src/github.com/cloudfoundry-community/kibana-me-logs
cf push kibana-myapp --no-start
cf bs kibana-myapp my-logstash-service
cf start kibana-myapp
```

Now view your Kibana UI in your browser.
