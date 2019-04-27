FROM openjdk:8-alpine
RUN apk add --no-cache curl bash nano jq && \
    cd /tmp && \
    curl http://mirrors.estointernet.in/apache/kafka/2.2.0/kafka_2.12-2.2.0.tgz -o kafka && \
    mkdir /kafka && \
    tar -xzf kafka --strip-components=1 -C /kafka/ && \
    rm kafka
COPY /build/package/shared/ /
WORKDIR /kafka
COPY /build/package/dev/kafka/start.sh \
/build/package/dev/kafka/config.sh \
/build/package/dev/kafka/broker-connected.sh \
./
ENTRYPOINT ["./start.sh"]