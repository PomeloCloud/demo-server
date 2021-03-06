FROM golang:1.9.0

WORKDIR /go/src/app

COPY . .

RUN go get github.com/astaxie/beego
RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 4561 8080

CMD ["sh", "run-after-sleep.sh"]