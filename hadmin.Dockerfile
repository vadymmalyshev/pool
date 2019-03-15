 FROM golang
 RUN mkdir -p /hadmin/config
 WORKDIR /hadmin
 COPY --from=pool-build-deps /pool/hadmin .
 COPY --from=pool-build-deps /pool/config/. config
 ENV hiveon-service=hadmin
 EXPOSE 3002
 CMD ["./hadmin"]