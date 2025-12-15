FROM ubuntu:24.04 AS shear

RUN apt-get update \
    && DEBIAN_FRONTEND="noninteractive" apt-get -y upgrade \
    && DEBIAN_FRONTEND="noninteractive" apt-get install -y build-essential wget git apt-transport-https pkg-config nano vim

# install golang
RUN wget -O /tmp/golang.tgz https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
RUN cd /usr/local && tar xzvf /tmp/golang.tgz
RUN update-alternatives --install "/usr/bin/go" "go" "/usr/local/go/bin/go" 0 \
    && update-alternatives --set go /usr/local/go/bin/go

WORKDIR /app
COPY . ./

RUN CGO_ENABLED=0 go build -o /app/shear *.go