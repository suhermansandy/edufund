FROM --platform=linux/amd64 golang AS build

WORKDIR /src

ENV CGO_ENABLED=0

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /out/app main.go

FROM scratch AS bin-unix

COPY --from=build /out/app /

FROM bin-unix AS bin-linux

FROM bin-linux AS bin

ENTRYPOINT ["/app"]