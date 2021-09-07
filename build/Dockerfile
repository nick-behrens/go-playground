FROM 888804278608.dkr.ecr.us-west-2.amazonaws.com/public/golang:1.16.1-buster AS BUILDER

ARG BUILD_NUMBER
ARG GIT_COMMIT
ARG GITHUB_TOKEN
ARG GOPRIVATE=github.com/snapdocs

COPY . /app
WORKDIR /app

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/snapdocs".insteadOf "https://github.com/snapdocs"
RUN mkdir -p bin/; go build -o bin/ -mod=readonly -ldflags="-X 'pkg.cmd.BuildNumber=${BUILD_NUMBER}' -X 'pkg.cmd.GitCommit=${GIT_COMMIT}'" ./...

FROM 888804278608.dkr.ecr.us-west-2.amazonaws.com/debian:bullseye
COPY --from=BUILDER /etc/ssl/certs/. /etc/ssl/certs/
COPY --from=BUILDER /app/bin/. /usr/local/bin/
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/go-playground" ]
