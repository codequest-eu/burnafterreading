FROM nanoservice/go:latest

# Create app directory.
ENV APP_HOME=/app
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME
ADD . $APP_HOME

# Fetch dependencies and build the server binary.
RUN go get -u github.com/codequest-eu/burnafterreading
RUN go get -u github.com/mitchellh/goamz/aws
RUN go get -u github.com/mitchellh/goamz/s3
RUN go build

# Get rid of build dependencies to keep the container size small.
RUN apk del --purge go
RUN apk del --purge git
RUN rm -rf /go
