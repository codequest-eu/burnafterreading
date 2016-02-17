FROM golang:latest

# Create app directory.
ENV APP_HOME=/app
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME
ADD . $APP_HOME

# Build dependencies and the ws binary
RUN go get -u github.com/codequest-eu/burnafterreading
RUN go get -u github.com/mitchellh/goamz/aws
RUN go get -u github.com/mitchellh/goamz/s3
RUN go build
