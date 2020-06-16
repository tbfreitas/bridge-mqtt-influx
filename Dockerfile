FROM vernemq/vernemq

ENV DOCKER_VERNEMQ_ACCEPT_EULA="yes"
ENV DOCKER_VERNEMQ_ALLOW_ANONYMOUS="on"

RUN rm -rf /etc/vernemq/vernemq.conf
COPY ./vernemq.conf /etc/vernemq/vernemq.conf
HEALTHCHECK CMD vernemq ping | grep -q pong

USER root

RUN apt-get update -y -q && apt-get upgrade -y -q 
RUN DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends apt-utils -y -q curl build-essential ca-certificates git 
RUN curl -s https://storage.googleapis.com/golang/go1.14.4.linux-amd64.tar.gz| tar -v -C /usr/local -xz

WORKDIR /vernemq

RUN mkdir /go && mkdir /go/src && mkdir /go/src/app

ENV GOPATH /go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
