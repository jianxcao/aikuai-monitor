FROM jakes/base-image:alpine-3.14.0-tz
LABEL maintainers="jianxcao"
LABEL description="iKuai monitor"
ARG TARGETPLATFORM
ENV GIN_MODE=release

ADD ./output /output
RUN cp -r /output/"$TARGETPLATFORM"/app /app

EXPOSE 7575

RUN chmod +x /app
CMD ["/app"]

