# go-gin-webapp
Go Gin Webapp template project

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
