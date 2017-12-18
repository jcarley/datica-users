FROM golang:1.9.2-stretch

ENV DEBIAN_FRONTEND noninteractive
ENV TERM xterm

ADD bashrc /root/.bashrc

RUN apt-get -qqy update && \
    apt-get install -qqy --no-install-recommends \
      ca-certificates                            \
      curl                                       \
      dnsutils                                   \
      git                                        \
      jq                                         \
      libgsf-1-dev                               \
      lsb-release                                \
      lsof                                       \
      net-tools                                  \
      netcat                                     \
      procps                                     \
      traceroute                                 \
      unzip                                      \
      vim                                        \
      wget                                       \
      postgresql-client-9.6                      \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN go get github.com/kardianos/govendor
