FROM google/golang

MAINTAINER Tyrell Keene <tyrell.wkeene@gmail.com>

RUN useradd -ms /bin/bash autobd
USER autobd

WORKDIR /home/autobd

ENV GOPATH=/home/autobd/go

RUN go get github.com/tywkeene/autobd

WORKDIR $GOPATH/src/github.com/tywkeene/autobd/
RUN bash build.sh

RUN mkdir /home/autobd/data
VOLUME /home/autobd/data

EXPOSE 8081
ENTRYPOINT ./autobd -root /home/autobd/data -api-port 8080
