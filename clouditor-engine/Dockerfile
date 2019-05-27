FROM openjdk:11-jre-slim

MAINTAINER Christian Banse <christian.banse@aisec.fraunhofer.de>

RUN apt update && apt -y --no-install-recommends install \
python3 \
python3-pip \
python3-dev \
gcc \
musl-dev \
libffi-dev \
libssl-dev \
nmap

RUN pip3 install --upgrade setuptools && pip3 install --upgrade sslyze nassl wheel

EXPOSE 9999

WORKDIR /usr/local/clouditor/

ADD build/distributions/engine-*.tar .
RUN mv engine-*/* . && rm -rf engine-*

ENV PATH /usr/local/clouditor/python:$PATH

ENTRYPOINT ["/usr/local/clouditor/bin/engine"]
