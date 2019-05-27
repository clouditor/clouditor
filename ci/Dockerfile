FROM java:8

RUN curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg | apt-key add -

RUN apt-get update && apt-get install -y \
                           apt-transport-https \
                           ca-certificates \
                           curl \
                           gnupg2 \
                           software-properties-common

RUN add-apt-repository \
       "deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") \
       $(lsb_release -cs) \
       stable"

RUN apt-get update && apt-get install -y nmap python3-pip docker-ce
RUN pip3 install docker-compose
