 FROM golang
 RUN mkdir -p /hapi/config /hapi/internal
 WORKDIR /hapi
 RUN apt update && \
      apt install -y libsasl2-dev libsasl2-modules libssl-dev && \
      git clone https://github.com/edenhill/librdkafka.git && \
      cd librdkafka && \
      ./configure --prefix /usr && \
      make && \
      make install
 COPY --from=pool-build-deps /pool/hapi .
 COPY ./config/. config/.
 COPY ./internal/. ./internal/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hapi
 EXPOSE 8080 8090
 CMD ["./hapi"]
