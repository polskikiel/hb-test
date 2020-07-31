FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add curl

RUN wget https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz -O - | tar -xzO linux-amd64/helm > /usr/local/bin/helm &&\
    chmod +x /usr/local/bin/helm

COPY ./helm-broker /helm-broker
COPY ./testing/ /testing/

CMD ["/helm-broker"]
