FROM openjdk:8-alpine
RUN apk add --no-cache curl bash nano jq && \
    cd /tmp && \
    curl http://mirrors.estointernet.in/apache/kafka/2.2.0/kafka_2.12-2.2.0.tgz -o kafka && \
    mkdir /kafka && \
    tar -xzf kafka --strip-components=1 -C /kafka/ && \
    rm kafka
WORKDIR /kafka
COPY ./build/package/dev/kafka/init.sh \
./build/package/dev/kafka/get-env.sh \
./build/package/dev/kafka/healthcheck.sh \
./build/package/dev/kafka/config.sh \
./
CMD ["./init.sh"]