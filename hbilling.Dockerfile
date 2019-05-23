 FROM golang
 RUN mkdir -p /hbilling/{config,internal}
 RUN apt update && \
     apt install -y libsasl2-dev libsasl2-modules libssl-dev && \
     git clone https://github.com/edenhill/librdkafka.git && \
     cd librdkafka && \
     ./configure --prefix /usr && \
     make && \
     make install && \
     cd ..
 WORKDIR /hbilling
 COPY --from=pool-build-deps /pool/hbilling .
 COPY ./config/. config/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hbilling
 CMD ["./hbilling"]