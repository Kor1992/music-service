FROM alpine:latest
WORKDIR /
COPY --from=builder /app/server /server
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/public /public
EXPOSE 8080
CMD ["/server"]