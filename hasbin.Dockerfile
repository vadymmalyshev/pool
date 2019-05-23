 FROM golang
 RUN mkdir -p /hasbin/config
 RUN apt update && \
    apt install -y libsasl2-dev libsasl2-modules libssl-dev && \
    git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install && \
    cd ..
 WORKDIR /hasbin
 COPY --from=pool-build-deps /pool/hasbin .
 COPY ./config/. config/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hasbin
 CMD ["./hasbin"]