# go-gin-webapp
Go Gin Webapp is a template project, that you can use to build your
RESTful microservices using [Gin](https://github.com/gin-gonic/gin) in
similar to [this blogpost](http://txt.fliglio.com/2014/07/restful-microservices-in-go-with-gin/).

It defaults the configuration to use https, if you configure the paths
to Certificate and Key correctly.

[Gin-gomonitor](https://github.com/zalando-techmonkeys/gin-gomonitor)
provides default metrics and you can easily implement your own
metrics, if you need.

It uses [gin-glog](https://github.com/zalando-techmonkeys/gin-glog)
middleware, which internally uses [glog](https://github.com/golang/glog) as
logger. This provides you leveled Logs, commandline options and logs
are flushed to disk in a separate goroutine with a configurable interval.

An [OAuth2 middleware](https://github.com/zalando-techmonkeys/gin-oauth2)
supports you to secure your web applications, if you need
authentication and authorization.

## Usage

To create a new web application, you can do the following steps

    APP=<appname>
    DST=$GOPATH/src/<repo-hoster>/<user>/$APP
    mkdir -p $DST
    go get -u github.com/szuecs/go-gin-webapp
    rsync -a --exclude=.git $GOPATH/src/github.com/szuecs/go-gin-webapp/ $DST
    cd $DST
    grep -rl go-gin-webapp | xargs sed -i "s@go-gin-webapp@$APP@g"
    mv cmd/go-gin-webapp cmd/$APP


The main package and function of the service is in
cmd/$APP/server.go. It parses flags and merges the configuration to
start the service.

The api package bootstraps the http(s) microservice, middleware,
registers the path specific handlers and you implement handlers
(similar to Controller in MVC) there.
