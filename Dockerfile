FROM registry.tor.ph/hiveon/pool/pool-base
RUN mkdir -p /pool
WORKDIR /pool
COPY . .
RUN eval `(ssh-agent)` && \
    ssh-add ~/.ssh/hiveon_ci_rsa && \
    git config --global url."git@git.tor.ph:".insteadOf "https://git.tor.ph/" && \ 
    go mod vendor && \ 
    go mod tidy && \
    go build cmd/hadmin/hadmin.go && \
    go build cmd/hasbin/hasbin.go && \
    go build cmd/hapi/hapi.go && \
    cd cmd/hbilling && \
    go build -o hbilling && \
    mv hbilling /pool/ && \
    ls /pool