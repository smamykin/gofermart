FROM ubuntu:20.04

COPY ./cmd/accrual/accrual_linux_amd64 /app/accrual/accrual_linux_amd64

ENV RUN_ADDRESS=0.0.0.0:8080

EXPOSE 8080

WORKDIR /app


CMD ["/app/accrual/accrual_linux_amd64", "-a", "$RUN_ADDRESS"]
