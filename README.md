Kibana Me Logs
==============

Draining logs from your Cloud Foundry hosted application to a backend Logstash/Elastic Search is great - it drains them all day long. Except, you can't see them. They are stored in Elastic Search and you have no way to access either Elastic Search nor a Kibana UI to view the logs.

**Events over time:**

![events-over-time](http://cl.ly/image/0r0O2a1n2D1W/events-over-time.png)

The application hosts Kibana 3 and a proxy that binds to your Logstash/Elastic Search backend service. You can now see your logs in Kibana!

This is a work-in-progress and/or a stop-gap until a better multi-tenant solution exists. It's primary weakness is it lacks authentication.

Security weaknesses
-------------------

If a user finds your Kibana app then they can see your application's logs without requiring username/password.

It requires the running of a shared proxy that grants users access to any backend service. Currently it doesn't have any authentication pass-thru.

But, it might provide you something useful in the meantime. Users can see and search their logs. And each other's logs.

Usage
-----

```
go get github.com/cloudfoundry-community/kibana-me-logs
cd $GOPATH/src/github.com/cloudfoundry-community/kibana-me-logs
cf push kibana-myapp --no-start
cf set-env kibana-myapp ES_PROXY <proxy-url>
cf bs kibana-myapp my-logstash-service
cf start kibana-myapp
```

Now view your Kibana UI in your browser. It should redirect to a url like `http://kibana-myapp.apps.1.2.3.4.xip.io/#/dashboard/file/logstash.json` automatically and start showing your logs.

See below for `<proxy-url>`.

Requirements
------------

This app currently requests a separate proxy app that Kibana will use to access Elastic Search.

```
go get github.com/longnguyen11288/proxy
cd $GOPATH/src/github.com/longnguyen11288/proxy
cf push proxy
```

The terminal will show the proxy URL, such as `proxy.apps.1.2.3.4.xip.io`

The proxy app will be shared by all users. See "security weaknesses" above.

Debugging
---------

-	"My Kibana dashboard is blank." It is possible you are not running the proxy above; or you haven't set the `$ES_PROXY` variable to the proxy hostname. See "Usage" above.
