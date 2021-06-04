FROM openjdk:13-slim

MAINTAINER Christian Banse <christian.banse@aisec.fraunhofer.de>

EXPOSE 9999

WORKDIR /usr/local/clouditor/

ADD build/distributions/engine-*.tar .
RUN mv engine-*/* . && rm -rf engine-*

ENTRYPOINT ["/usr/local/clouditor/bin/engine"]
