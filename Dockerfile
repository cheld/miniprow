FROM golang:1.14.3-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/cicd-bot .

FROM scratch AS bin
COPY --from=build /out/cicd-bot /
ENTRYPOINT [ "/cicd-bot" ]
CMD [ "serve"]

FROM bin AS custom
COPY ./cicd-bot.yaml /etc

