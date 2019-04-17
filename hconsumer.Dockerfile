 FROM golang
 RUN mkdir -p /hconsumer/config /hconsumer/internal
 WORKDIR /hconsumer
 COPY --from=pool-build-deps /pool/hconsumer .
 COPY ./config/. config/.
 COPY ./internal/. ./internal/.
 ENV hiveon-service=hconsumer
 CMD ["./hconsumer"]
