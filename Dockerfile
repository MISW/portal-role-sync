FROM golang:1.18 AS builder

COPY . /work
WORKDIR /work
ENV CGO_ENABLED=0
RUN go build -o /portal-role-sync

FROM gcr.io/distroless/static:debug

COPY --from=builder /portal-role-sync /bin

ENTRYPOINT [ "/bin/portal-role-sync" ]
