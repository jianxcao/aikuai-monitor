FROM alpine
LABEL maintainers="jianxcao"
LABEL description="iKuai monitor"
ARG TARGETPLATFORM
ENV GIN_MODE=release
ENV AIKUAI_MONITOR_CONFIG_PATH=/app/config

RUN apk --no-cache add ca-certificates

ADD ./output /output

RUN cp -r /output/"$TARGETPLATFORM" /app

EXPOSE 7575
WORKDIR /app

RUN chmod +x /app/app

CMD ["/app/app"]

