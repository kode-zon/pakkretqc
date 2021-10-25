FROM golang:1.17

WORKDIR /app
COPY pakkretqc-master /app/pakkretqc

ARG PAKKRETQC_ALM_ENDPOINT=https://yourhost:9999/qcweb/
WORKDIR /app/pakkretqc
RUN . ./set_proxy.sh
RUN ./build.sh

WORKDIR /app/pakkretqc/dist/pakkretqc
RUN cd bin && ./pakkretqc