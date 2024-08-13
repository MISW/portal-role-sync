# workspace
FROM golang:1.23 AS workspace

COPY . /portal-role-sync

WORKDIR /portal-role-sync

RUN go mod download \
  && CGO_ENABLED=0 go build -buildmode pie -o /portal-role-sync/portal-role-sync

# production
FROM gcr.io/distroless/base:debug AS production

RUN ["/busybox/sh", "-c", "ln -s /busybox/sh /bin/sh"]
RUN ["/busybox/sh", "-c", "ln -s /bin/env /usr/bin/env"]

COPY --from=workspace /portal-role-sync/portal-role-sync /bin/portal-role-sync

ENTRYPOINT ["/bin/portal-role-sync"]
