 FROM golang
 RUN mkdir -p /hadmin/config
 WORKDIR /hadmin
 COPY --from=pool-build-deps /pool/hadmin .
 COPY ./config/. config/.
 RUN mv ./config/config.dev.yaml ./config/config.yaml
 ENV hiveon-service=hadmin
 EXPOSE 3002
 CMD ["./hadmin"]