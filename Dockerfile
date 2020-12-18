FROM ubuntu
COPY ./bin/miniprow /
ENTRYPOINT [ "/miniprow" ]
CMD [ "serve"]

