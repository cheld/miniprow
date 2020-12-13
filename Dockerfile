FROM scratch
COPY ./bin/miniprow /
ENTRYPOINT [ "/miniprow" ]
CMD [ "serve"]

