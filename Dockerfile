FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY ./helm-broker /helm-broker
COPY ./testing/ /testing/

CMD ["/helm-broker"]
