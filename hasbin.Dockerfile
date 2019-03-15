 FROM golang
 RUN mkdir -p /hasbin/config
 WORKDIR /hasbin
 COPY --from=pool-build-deps /pool/hasbin .
 COPY --from=pool-build-deps /pool/config/. config
 ENV hiveon-service=hasbin
 CMD ["./hasbin"]