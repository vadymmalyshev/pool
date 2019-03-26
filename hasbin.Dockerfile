 FROM golang
 RUN mkdir -p /hasbin/config
 WORKDIR /hasbin
 COPY --from=pool-build-deps /pool/hasbin .
 COPY ./config/. config/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hasbin
 CMD ["./hasbin"]