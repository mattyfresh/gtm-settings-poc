FROM golang:1.11

ARG GITHUB_ACCESS_TOKEN
ARG GOOGLE_APPLICATION_CREDENTIALS
ARG GTM_ACCOUNT_ID
ARG SLACK_BOT_API_TOKEN

ENV GITHUB_ACCESS_TOKEN $GITHUB_ACCESS_TOKEN
ENV GOOGLE_APPLICATION_CREDENTIALS $GOOGLE_APPLICATION_CREDENTIALS
ENV GTM_ACCOUNT_ID $GTM_ACCOUNT_ID
ENV SLACK_BOT_API_TOKEN $SLACK_BOT_API_TOKEN

RUN mkdir go/src/app

ADD . go/src/app/

WORKDIR /go/src/app

RUN go get ./...

RUN go build -o gobot main.go gtm_controller.go gtm_validators.go gtm_service.go

CMD ["/app/gobot"]

