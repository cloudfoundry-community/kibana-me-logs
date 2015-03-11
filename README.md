Kibana Me Logs
==============

Draining logs from your Cloud Foundry hosted application to a backend Logstash/Elastic Search is great - it drains them all day long. Except, you can't see them. They are stored in Elastic Search and you have no way to access either Elastic Search nor a Kibana UI to view the logs.

**Events over time:**

![events-over-time](http://cl.ly/image/0r0O2a1n2D1W/events-over-time.png)

**Line-by-line logs**

![line-by-line](http://cl.ly/image/2k0K3t0g1V0V/line-by-line_logs.png)

The application hosts Kibana 3 and a proxy that binds to your Logstash/Elastic Search backend service. You can now see your logs in Kibana!

This is a work-in-progress and/or a stop-gap until a better multi-tenant solution exists. It's primary weakness is it lacks authentication.

It assumes that users are getting Logstash via the Docker/Logstash Service Broker. See below for details.

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
cf bs kibana-myapp my-logstash-service
cf start kibana-myapp
```

Now view your Kibana UI in your browser. It should redirect to a url like `http://kibana-myapp.apps.1.2.3.4.xip.io/#/dashboard/file/logstash.json` automatically and start showing your logs.

Docker/Logstash Service Broker
------------------------------

A requirement for this application is that the service binding credentials to logstash/elasticsearch fit a certain schema. This is the schema that comes from the Docker/Logstash Service Broker.

The easiest way to deploy this service broker is with the docker-services-boshworkspace.

See its README for detailed instructions.

When running `bosh setup deployment` choose "Logstash 1.4" as the service to be deployed.

Logstash Service Brokers
------------------------

If you are building your own service broker for Logstash/Elastic Search, then this application assumes that the binding credentials look like:

```json
{
  "hostname": "10.10.5.251",
  "ports": {
   "514/tcp": "49160",
   "9200/tcp": "49161",
   "9300/tcp": "49159"
  }
}
```

The `9200/tcp` port is the port for the Elastic Search API.

If your Elastic Search API is hosted on a different hostname than the remainder of the ELK cluster then please submit a PR/issue so we can discuss what to do next. Happy to help.
