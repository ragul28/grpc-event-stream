FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app
USER app

ARG SERVICE
COPY ./bin/${SERVICE} /server

ENV PORT 8080
ENTRYPOINT ["/server"]