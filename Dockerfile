FROM registry.tor.ph/hiveon/pool/pool-base
RUN mkdir -p /pool
WORKDIR /pool
COPY . .
RUN eval `(ssh-agent)` && \
    ssh-add ~/.ssh/hiveon_ci_rsa && \
    git config --global url."git@git.tor.ph:".insteadOf "https://git.tor.ph/" && \ 
    apt update && \
    apt install -y libsasl2-dev libsasl2-modules libssl-dev && \
    git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install && \
    cd .. && \
    go mod vendor && \ 
    go mod tidy && \
    go build cmd/hadmin/hadmin.go && \
    go build cmd/hasbin/hasbin.go && \
    go build cmd/hapi/hapi.go && \
    go build cmd/hconsumer/hconsumer.go && \
    cd cmd/hbilling && \
    go build -o hbilling && \
    mv hbilling /pool/ && \
    ls /pool
