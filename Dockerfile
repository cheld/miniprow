FROM golang:1.14.3-alpine AS build
WORKDIR /src
COPY . .
RUN go build -v -o /out/miniprow cmd/miniprow/miniprow.go

FROM scratch
COPY --from=build /out/miniprow /
ENTRYPOINT [ "/miniprow" ]
CMD [ "serve"]

