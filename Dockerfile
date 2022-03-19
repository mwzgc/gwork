FROM alpine:3.13

WORKDIR /opt/

COPY bin/server /opt

RUN apk --no-cache add tini tzdata

ENV PATH /opt:$PATH

ENTRYPOINT ["/sbin/tini", "--", "/opt/server"]

CMD ["server"]