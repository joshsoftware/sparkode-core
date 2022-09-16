# Check for latest version here: https://hub.docker.com/_/buildpack-deps?tab=tags&page=1&name=buster&ordering=last_updated
# This is just a snapshot of buildpack-deps:buster that was last updated on 2019-12-28.
FROM judge0/buildpack-deps:buster-2019-12-28

RUN mkdir /app
WORKDIR /app
# Check for latest version here: https://gcc.gnu.org/releases.html, https://ftpmirror.gnu.org/gcc
ENV GCC_VERSIONS \
      7.4.0 \
      8.3.0 \
      9.2.0
RUN set -xe && \
    for VERSION in $GCC_VERSIONS; do \
      curl -fSsL "https://ftpmirror.gnu.org/gcc/gcc-$VERSION/gcc-$VERSION.tar.gz" -o /tmp/gcc-$VERSION.tar.gz && \
      mkdir /tmp/gcc-$VERSION && \
      tar -xf /tmp/gcc-$VERSION.tar.gz -C /tmp/gcc-$VERSION --strip-components=1 && \
      rm /tmp/gcc-$VERSION.tar.gz && \
      cd /tmp/gcc-$VERSION && \
      ./contrib/download_prerequisites && \
      { rm *.tar.* || true; } && \
      tmpdir="$(mktemp -d)" && \
      cd "$tmpdir"; \
      if [ $VERSION = "9.2.0" ]; then \
        ENABLE_FORTRAN=",fortran"; \
      else \
        ENABLE_FORTRAN=""; \
      fi; \
      /tmp/gcc-$VERSION/configure \
        --disable-multilib \
        --enable-languages=c,c++$ENABLE_FORTRAN \
        --prefix=/usr/local/gcc-$VERSION && \
      make -j$(nproc) && \
      make -j$(nproc) install-strip && \
      rm -rf /tmp/*; \
    done

# Check for latest version here: https://www.ruby-lang.org/en/downloads
ENV RUBY_VERSIONS \
      2.7.0
RUN set -xe && \
    for VERSION in $RUBY_VERSIONS; do \
      curl -fSsL "https://cache.ruby-lang.org/pub/ruby/${VERSION%.*}/ruby-$VERSION.tar.gz" -o /tmp/ruby-$VERSION.tar.gz && \
      mkdir /tmp/ruby-$VERSION && \
      tar -xf /tmp/ruby-$VERSION.tar.gz -C /tmp/ruby-$VERSION --strip-components=1 && \
      rm /tmp/ruby-$VERSION.tar.gz && \
      cd /tmp/ruby-$VERSION && \
      ./configure \
        --disable-install-doc \
        --prefix=/usr/local/ruby-$VERSION && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*; \
    done

# Check for latest version here: https://www.python.org/downloads
ENV PYTHON_VERSIONS \
      3.8.1 \
      2.7.17
RUN set -xe && \
    for VERSION in $PYTHON_VERSIONS; do \
      curl -fSsL "https://www.python.org/ftp/python/$VERSION/Python-$VERSION.tar.xz" -o /tmp/python-$VERSION.tar.xz && \
      mkdir /tmp/python-$VERSION && \
      tar -xf /tmp/python-$VERSION.tar.xz -C /tmp/python-$VERSION --strip-components=1 && \
      rm /tmp/python-$VERSION.tar.xz && \
      cd /tmp/python-$VERSION && \
      ./configure \
        --prefix=/usr/local/python-$VERSION && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*; \
    done


# Check for latest version here: https://jdk.java.net
RUN set -xe && \
    curl -fSsL "https://download.java.net/java/GA/jdk13.0.1/cec27d702aa74d5a8630c65ae61e4305/9/GPL/openjdk-13.0.1_linux-x64_bin.tar.gz" -o /tmp/openjdk13.tar.gz && \
    mkdir /usr/local/openjdk13 && \
    tar -xf /tmp/openjdk13.tar.gz -C /usr/local/openjdk13 --strip-components=1 && \
    rm /tmp/openjdk13.tar.gz && \
    ln -s /usr/local/openjdk13/bin/javac /usr/local/bin/javac && \
    ln -s /usr/local/openjdk13/bin/java /usr/local/bin/java && \
    ln -s /usr/local/openjdk13/bin/jar /usr/local/bin/jar

# Check for latest version here: https://nodejs.org/en
ENV NODE_VERSIONS \
      12.14.0
RUN set -xe && \
    for VERSION in $NODE_VERSIONS; do \
      curl -fSsL "https://nodejs.org/dist/v$VERSION/node-v$VERSION.tar.gz" -o /tmp/node-$VERSION.tar.gz && \
      mkdir /tmp/node-$VERSION && \
      tar -xf /tmp/node-$VERSION.tar.gz -C /tmp/node-$VERSION --strip-components=1 && \
      rm /tmp/node-$VERSION.tar.gz && \
      cd /tmp/node-$VERSION && \
      ./configure \
        --prefix=/usr/local/node-$VERSION && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*; \
    done

# Check for latest version here: https://golang.org/dl
ENV GO_VERSIONS \
      1.13.5
RUN set -xe && \
    for VERSION in $GO_VERSIONS; do \
      curl -fSsL "https://storage.googleapis.com/golang/go$VERSION.linux-amd64.tar.gz" -o /tmp/go-$VERSION.tar.gz && \
      mkdir /usr/local/go-$VERSION && \
      tar -xf /tmp/go-$VERSION.tar.gz -C /usr/local/go-$VERSION --strip-components=1 && \
      rm -rf /tmp/*; \
    done

# Check for latest version here: https://packages.debian.org/buster/clang-7
# Used for additional compilers for C, C++ and used for Objective-C.
RUN set -xe && \
    apt-get update && \
    apt-get install -y --no-install-recommends clang-7 gnustep-devel && \
    rm -rf /var/lib/apt/lists/*

RUN set -xe && \
    apt-get update && \
    apt-get install -y --no-install-recommends locales && \
    rm -rf /var/lib/apt/lists/* && \
    echo "en_US.UTF-8 UTF-8" > /etc/locale.gen && \
    locale-gen
ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8

RUN set -xe && \
    apt-get update && \
    apt-get install -y --no-install-recommends git libcap-dev && \
    rm -rf /var/lib/apt/lists/* && \
    git clone https://github.com/judge0/isolate.git /tmp/isolate && \
    cd /tmp/isolate && \
    git checkout ad39cc4d0fbb577fb545910095c9da5ef8fc9a1a && \
    make -j$(nproc) install && \
    rm -rf /tmp/*
ENV BOX_ROOT /var/local/lib/isolate

LABEL maintainer="Herman Zvonimir Došilović <hermanz.dosilovic@gmail.com>"
LABEL version="1.4.0"
