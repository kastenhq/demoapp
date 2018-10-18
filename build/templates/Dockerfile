FROM scratch
COPY ${BIN} /
EXPOSE 8080
ENTRYPOINT [ "/${BIN}" ]
CMD [ "--port", "8080", "--host", "0.0.0.0"]
