 FROM golang
 RUN mkdir -p /hadmin/config
 WORKDIR /hadmin
 RUN apt update && \
      apt install -y libsasl2-dev libsasl2-modules libssl-dev && \
      git clone https://github.com/edenhill/librdkafka.git && \
      cd librdkafka && \
      ./configure --prefix /usr && \
      make && \
      make install
 COPY --from=pool-build-deps /pool/hadmin .
 COPY ./config/. config/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hadmin
 EXPOSE 3002
 CMD ["./hadmin"]
