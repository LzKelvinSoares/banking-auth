FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server .

FROM gcr.io/distroless/static-debian12
COPY --from=build /server /server
# COPY --from=build /src/resources /resources   # só se a API ler arquivos de resources/
EXPOSE 8080
ENTRYPOINT ["/server"]