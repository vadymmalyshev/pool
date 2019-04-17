#hiveon-api image
FROM  golang as consumer-build-deps
RUN mkdir -p /hconsumer
WORKDIR /hconsumer
COPY . /hconsumer
RUN apt update && \
    apt install -y libsasl2-dev libsasl2-modules libssl-dev && \
    git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install && \
COPY --from=pool-build-deps /pool/hconsumer /hconsumer
COPY --from=pool-build-deps /pool/config/. /hconsumer/config/.
ENV hiveon-service=hbilling
CMD ["./hconsumer"]
