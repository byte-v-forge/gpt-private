FROM docker.m.daocloud.io/library/golang:1.26-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
ENV PATH=/root/go/bin:$PATH

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache git protobuf-dev \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

COPY common-lib /common-lib
COPY gpt/proto /gpt/proto
COPY gpt-private/gopay/go.mod gpt-private/gopay/go.sum ./
RUN go mod edit -replace github.com/byte-v-forge/common-lib=/common-lib \
    && go mod download

COPY gpt-private/gopay ./
RUN go mod edit -replace github.com/byte-v-forge/common-lib=/common-lib \
    && mkdir -p pb \
    && rm -f pb/*.pb.go pb/*_grpc.pb.go \
    && protoc -I /gpt/proto -I /common-lib/proto --go_out=pb --go-grpc_out=pb /gpt/proto/payment.proto \
    && go build -o /out/gpt-checkout ./cmd/gopay-payment-server

FROM docker.m.daocloud.io/library/alpine:latest AS runtime

WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates

COPY --from=builder /out/gpt-checkout /app/bin/gpt-checkout

EXPOSE 50054
CMD ["/app/bin/gpt-checkout"]
