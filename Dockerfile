FROM ubuntu:24.04 AS shear

RUN apt-get update \
    && DEBIAN_FRONTEND="noninteractive" apt-get -y upgrade \
    && DEBIAN_FRONTEND="noninteractive" apt-get install -y build-essential wget pkg-config

RUN wget -O /tmp/golang.tgz https://go.dev/dl/go1.23.0.linux-amd64.tar.gz \
    && cd /usr/local && tar xzvf /tmp/golang.tgz \
    && rm /tmp/golang.tgz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/shear ./bot/main.go

ENTRYPOINT ["/app/shear"]