 FROM golang
 RUN mkdir -p /hbilling/config
 WORKDIR /hbilling
 COPY --from=pool-build-deps /pool/hbilling .
 COPY --from=pool-build-deps /pool/config/. config
 ENV hiveon-service=hbilling
 CMD ["./hbilling"]