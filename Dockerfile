FROM golang:alpine
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /
COPY . /
RUN go build .

FROM alpine
ENV TZ=Asia/Shanghai \
    RAWURL='your rawurl' \
    SECRETID='your secretid' \
    SECRETKEY='your secretkey'
WORKDIR /app/bin/
COPY --from=0 /zu_logic /app/bin/
COPY ./config/config.toml /app/config/
EXPOSE 6900 6901
CMD ["./zu_logic"]
