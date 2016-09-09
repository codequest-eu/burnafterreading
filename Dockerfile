FROM nanoservice/go:latest

# Fetch dependencies
RUN go get -u github.com/mitchellh/goamz/aws
RUN go get -u github.com/mitchellh/goamz/s3

# Create app directory.
ENV CODE_HOME=$GOPATH/src/github.com/codequest-eu/burnafterreading
RUN mkdir -p $CODE_HOME
WORKDIR $CODE_HOME
ADD . $CODE_HOME

# Build the server binary.
RUN go build

# Create app directory.
ENV APP_HOME=/app
RUN mkdir $APP_HOME
RUN mv $CODE_HOME/burnafterreading $APP_HOME/app
WORKDIR $APP_HOME

# Get rid of build dependencies to keep the container size small.
RUN apk del --purge go
RUN apk del --purge git
RUN rm -rf /go
