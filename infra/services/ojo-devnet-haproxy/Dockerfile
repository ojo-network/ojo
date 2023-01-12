FROM amd64/haproxy:2.6-alpine

COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg
RUN haproxy -f /usr/local/etc/haproxy/haproxy.cfg -c
