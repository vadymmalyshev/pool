 FROM golang
 RUN mkdir -p /hconsumer/config
 WORKDIR /hbilling
 COPY --from=pool-build-deps /pool/hconsumer .
 COPY ./config/. config/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hconsumer
 CMD ["./hconsumer"]
