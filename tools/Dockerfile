FROM ubuntu:bionic

ENV GO_VERSION 1.13

RUN apt update && apt install -y libboost-dev \
    libboost-program-options-dev \
    libboost-system-dev \
    libboost-thread-dev \
    libboost-math-dev \
    libboost-test-dev \
    zlib1g-dev \
    cmake \
    g++ \
    git \
    pkg-config \
    wget \
    rapidjson-dev 

RUN git clone --branch "8.9.0-fixes" --recursive https://github.com/orktes/vowpal_wabbit.git && \
    cd vowpal_wabbit && make install && ldconfig && \
    cd .. && rm -Rf vowpal_wabbit

RUN  wget -qO- https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz \
    | tar -C /usr/local -xz && \
    export PATH="/usr/local/go/bin:$PATH" && \
    go version; 

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
