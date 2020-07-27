FROM registry.access.redhat.com/ubi8/ubi:latest

COPY bin/request-monitor /
WORKDIR /
