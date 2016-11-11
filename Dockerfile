FROM registry.opensource.zalan.do/stups/alpine:UPSTREAM

ADD config.yaml.sample /etc/go-gin-webapp/config.yaml

# add scm-source
ADD scm-source.json /

# add binary
ADD build/linux/go-gin-webapp /

ENTRYPOINT ["/go-gin-webapp"]
