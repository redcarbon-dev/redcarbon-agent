FROM golang:1.27

RUN useradd --create-home e2e
USER e2e
WORKDIR /home/e2e/src

COPY --chown=e2e go.mod go.sum ./
RUN go mod download

COPY --chown=e2e . .

CMD ["go", "test", "-tags", "e2e", "-count=1", "-v", "./test/e2e/..."]
