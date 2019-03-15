 FROM golang
 RUN mkdir -p /hapi/config
 WORKDIR /hapi
 COPY --from=pool-build-deps /pool/hapi .
 COPY --from=pool-build-deps /pool/config/. config
 ENV hiveon-service=hapi
 EXPOSE 8080 8090
 CMD ["./hapi"]
