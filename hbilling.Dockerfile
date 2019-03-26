 FROM golang
 RUN mkdir -p /hbilling/config
 WORKDIR /hbilling
 COPY --from=pool-build-deps /pool/hbilling .
 COPY ./config/. config/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hbilling
 CMD ["./hbilling"]