FROM alpine:latest
ARG VERSION=latest
RUN apk add \
  curl \
  bash \
  ca-certificates \
  openssl \
  postgresql-client \
  mysql-client \
  redis \
  mongodb-tools \
  sqlite \
  # replace busybox utils
  tar \
  gzip \
  pigz \
  bzip2 \
  # there is no pbzip2 yet
  lzip \
  xz-dev \
  lzop \
  xz \
  # pixz is in edge atm
  zstd \
  # microsoft sql dependencies \
  libstdc++ \
  gcompat \
  icu \
  && \
  rm -rf /var/cache/apk/*

WORKDIR /tmp
RUN wget https://aka.ms/sqlpackage-linux && \
    unzip sqlpackage-linux -d /opt/sqlpackage && \
    rm sqlpackage-linux && \
    chmod +x /opt/sqlpackage/sqlpackage \
    curl -sSL https://d.juicefs.com/install | sh -

ENV PATH="${PATH}:/opt/sqlpackage"

ADD install /install
RUN /install ${VERSION} && rm /install

# CMD ["/usr/local/bin/gobackup"]