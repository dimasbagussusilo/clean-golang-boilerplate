FROM 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/golang-alpine as builder

ARG REPOSITORY_NAME

WORKDIR /$REPOSITORY_NAME/
COPY . /$REPOSITORY_NAME/

RUN protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. files/grpc-protos/*.proto
RUN wire ./...

RUN go build -o /$REPOSITORY_NAME/$REPOSITORY_NAME ./app/cmd/

FROM 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/alpine

ARG APP_ENV
ARG REPOSITORY_NAME

WORKDIR /srv/app

RUN mkdir -pv /src/files/ssl
RUN mkdir -pv /src/files/templates

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /$REPOSITORY_NAME/$REPOSITORY_NAME/ /$REPOSITORY_NAME
COPY --from=builder /$REPOSITORY_NAME/.env.$APP_ENV /.env

EXPOSE 8000 8001

ENTRYPOINT ["$ENT"]