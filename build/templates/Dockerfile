FROM scratch
COPY ${BIN} /
EXPOSE 8000
ENTRYPOINT [ "/${BIN}" ]
CMD [ "--port", "8000", "--host", "0.0.0.0"]
