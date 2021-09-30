FROM golang:1.12.7 as builder
RUN mkdir -p /work
WORKDIR /work

COPY go.mod go.sum Makefile /work/
RUN make fetchdeps

COPY . /work/
RUN make clean build


FROM scratch as runner
WORKDIR /

# If you require hard-coded vault or k8s CAs, place them in `main/resources/ca-bundle.pem` and uncomment this line.
# COPY main/resources/ca-bundle.pem /etc/ssl/ca-bundle.pem

COPY --from=builder /work/build/certgen-static /certgen-static
ENTRYPOINT [ "/certgen-static" ]
