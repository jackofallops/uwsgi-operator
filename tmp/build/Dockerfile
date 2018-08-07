FROM alpine:3.6

RUN adduser -D uwsgi-operator
USER uwsgi-operator

ADD tmp/_output/bin/uwsgi-operator /usr/local/bin/uwsgi-operator
