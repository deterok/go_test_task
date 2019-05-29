FROM golang

RUN mkdir -p /go/src/github.com/deterok/go_test_task

ADD . /go/src/github.com/deterok/go_test_task/
WORKDIR /go/src/github.com/deterok/go_test_task

RUN go get -u github.com/golang/dep/...
RUN dep ensure

# We use wathcer to monitor file changes. If changes occur, we rebuild the project.
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

CMD  ["watcher", "-run", "github.com/deterok/go_test_task/payments/cmd", "-watch", "github.com/deterok/go_test_task/payments"]
