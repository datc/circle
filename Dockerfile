FROM golang

# Build app
RUN mkdir -p /usr/static/app/bk
ENV GOPATH /usr/static/app
WORKDIR /usr/static/app/bk

RUN git clone --depth 1 git://github.com/datc/circle.git . && go get github.com/datc/circle && go build -o circle

EXPOSE 80

CMD ["/usr/static/app/bk/circle"]