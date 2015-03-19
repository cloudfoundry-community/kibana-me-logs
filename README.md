Kibana Me Logs
==============

Draining logs from your Cloud Foundry hosted application to a backend Logstash/Elastic Search is great - it drains them all day long. Except, you can't see them. They are stored in Elastic Search and you have no way to access either Elastic Search nor a Kibana UI to view the logs.

**Events over time:**

![events-over-time](http://cl.ly/image/0r0O2a1n2D1W/events-over-time.png)

**Line-by-line logs**

![line-by-line](http://cl.ly/image/3M2U3A3u1v0S/line-by-line_logs.png)

The application hosts Kibana 3 and a proxy that binds to your Logstash/Elastic Search backend service. You can now see your logs in Kibana!

This is a work-in-progress and/or a stop-gap until a better multi-tenant solution exists. It's primary weakness is it lacks authentication.

It assumes that users are getting Logstash via the Docker/Logstash Service Broker. See below for details.

For users, there is the complimentary `cf kibana-me-logs` CLI plugin https://github.com/cloudfoundry-community/cf-plugin-kibana-me-logs

For administrators, there is `./bin/upgrade-all.sh` to systematically upgrade all kibana-me-logs applications on your Cloud Foundry.

View logs for a specific app
----------------------------

If you are binding one logstash service to many applications in a space, then you might want a way to only see the logs for one app at a time.

`http://kibana.DOMAIN/#/dashboard/file/app-logs-01a4ad6a-51b1-450b-ab8d-ef5b836bb8cb.json`

Put the GUID for the application into the URL above.

Security weaknesses
-------------------

If a user finds your Kibana app then they can see your application's logs without requiring username/password.

It requires the running of a shared proxy that grants users access to any backend service. Currently it doesn't have any authentication pass-thru.

But, it might provide you something useful in the meantime. Users can see and search their logs. And each other's logs.

Assumptions
-----------

It is assumed that:

-	you already have an application running on your own Cloud Foundry
-	you are using the Docker/Logstash Service Broker (see below) or a service broker that provides matching credentials schema
-	you have already bound your application to your logstash service instance and Cloud Foundry is already draining logs into it
-	you have the Go programming language installed on your machine

For example:

```
cf create-service logstash14 free my-logstash-service
cf bind-service my-app my-logstash-service
cf restart my-app
```

You can confirm that your application is bound to the service:

```
$ cf services
Getting services in org system / space dev as admin...
OK

name                  service      plan   bound apps   status
my-logstash-service   logstash14   free   my-app       available
```

Usage
-----

To view your application's logs in Kibana you need to deploy the `kibana-me-logs` application and also bind it to the same `my-logstash-service` service instance as above:

```
cd /tmp; rm -rf kibana-me-logs
git clone https://github.com/cloudfoundry-community/kibana-me-logs
cd kibana-me-logs
cf push kibana-myapp --no-start --random-route -b https://github.com/cloudfoundry/go-buildpack
cf bs kibana-myapp my-logstash-service
cf start kibana-myapp
```

Now view your Kibana UI in your browser. It should redirect to a url like `http://kibana-myapp.apps.1.2.3.4.xip.io/#/dashboard/file/logstash.json` automatically and start showing your logs.

If you are a regular Go user, you can also fetch the application using:

```
go get github.com/cloudfoundry-community/kibana-me-logs
cd $GOPATH/src/github.com/cloudfoundry-community/kibana-me-logs
```

Docker/Logstash Service Broker
------------------------------

A requirement for this application is that the service binding credentials to logstash/elasticsearch fit a certain schema. This is the schema that comes from the Docker/Logstash Service Broker.

The easiest way to deploy this service broker is with the docker-services-boshworkspace.

See its README for detailed setup instructions for administrators.

When running `bosh setup deployment` choose "Logstash 1.4" as the service to be deployed.

Logstash Service Brokers
------------------------

If you are building your own service broker for Logstash/Elastic Search, then this application assumes that the binding credentials look like:

```json
{
  "name": "my-logstash-service",
  "credentials": {
    "hostname": "10.10.5.251",
    "ports": {
     "514/tcp": "49160",
     "9200/tcp": "49161",
     "9300/tcp": "49159"
    }
  },
  "syslog_drain_url": "http://10.10.5.251:49160",
  ...
}
```

The `9200/tcp` port is the port for the Elastic Search API.

If your Elastic Search API is hosted on a different hostname than the remainder of the ELK cluster then please submit a PR/issue so we can discuss what to do next. Happy to help.

Your service bindings should also include a `syslog_drain_url` URI. Cloud Foundry will use this to automatically setup a continuous syslog drain of your application's logs into your logstash/elastic search service instance.
