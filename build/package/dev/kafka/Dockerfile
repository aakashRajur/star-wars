FROM openjdk:8-alpine
RUN apk add --no-cache curl bash nano jq perl-getopt-long && \
    cd /tmp && \
    curl https://www-eu.apache.org/dist/kafka/2.2.0/kafka_2.12-2.2.0.tgz -o kafka && \
    mkdir /kafka && \
    tar -xzf kafka --strip-components=1 -C /kafka/ && \
    rm kafka
COPY /build/package/shared/ /util/
WORKDIR /kafka
COPY /build/package/dev/kafka/*.sh ./
ENTRYPOINT ["./start.sh"]