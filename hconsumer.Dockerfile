 FROM golang
 RUN mkdir -p /hconsumer/config /hconsumer/internal
 WORKDIR /hconsumer
 RUN apt update && \
     apt install -y libsasl2-dev libsasl2-modules libssl-dev && \
     git clone https://github.com/edenhill/librdkafka.git && \
     cd librdkafka && \
     ./configure --prefix /usr && \
     make && \
     make install && \
     cd .. 
 COPY --from=pool-build-deps /pool/hconsumer .
 COPY ./config/. config/.
 COPY ./internal/. ./internal/.
 ENV hiveon-service=hconsumer
 CMD ["./hconsumer"]
