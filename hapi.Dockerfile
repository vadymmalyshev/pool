 FROM golang
 RUN mkdir -p /hapi/config /hapi/internal
 WORKDIR /hapi
 COPY --from=pool-build-deps /pool/hapi .
 COPY ./config/. config/.
 COPY ./internal/. ./internal/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hapi
 EXPOSE 8080 8090
 CMD ["./hapi"]
