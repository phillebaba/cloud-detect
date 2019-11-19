FROM golang:1.13-alpine AS build

COPY . /go/src/github.com/phillebaba/cloud-detect/
WORKDIR /go/src/github.com/phillebaba/cloud-detect/
RUN CGO_ENABLED=0 go build -o /bin/cloud-detect main.go

FROM scratch
COPY --from=build /bin/cloud-detect /bin/cloud-detect

ENTRYPOINT ["/bin/cloud-detect"]
