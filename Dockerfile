FROM golang:onbuild

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o gobot main.go gtm_controller.go gtm_validators.go gtm_service.go

CMD ["/app/gobot"]

