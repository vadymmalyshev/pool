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
    cd /hconsumer && \
    go mod vendor && \
    go mod tidy && \
    go build hconsumer && \
    ls && \
    cp kafka/* conf/ && \
    cp -r kafka conf/ 

CMD ["./consumer"]
#FROM golang
#RUN mkdir -p /consumer/conf
#WORKDIR /consumer
#COPY --from=consumer-build-deps /consumer/conf conf/.
#COPY --from=consumer-build-deps /consumer/consumer .
#RUN mv ./conf/config.dev.yaml ./conf/config.yaml
#COPY --from=consumer-build-deps /consumer/kafka ./conf/kafka
#COPY --from=consumer-build-deps /consumer/kafka ./kafka
#ENV build-number=${CI_PIPELINE_ID:-latest}
#CMD ["./consumer"]
