FROM --platform=linux/amd64 ubuntu:20.04
COPY ./cmd/accrual/accrual_linux_amd64 .

CMD ["./accrual_linux_amd64"]