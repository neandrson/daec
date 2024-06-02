FROM quay.io/libpod/alpine

RUN apk add build-base
RUN wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /go
#ENV PATH $PATH:/go/bin

COPY . /daec
WORKDIR /daec

#RUN go build -o orchestrator cmd/orchestrator/main.go

#CMD source scripts/setenv.sh && ./orchestrator
CMD source scripts/setenv.sh && go run -race cmd/orchestrator/main.go
