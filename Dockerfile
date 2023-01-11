# syntax=docker/dockerfile:1.3

FROM alpine:3.14
EXPOSE 8080
EXPOSE 8081

COPY module-runner .
RUN chmod -R 775 module-runner
ENTRYPOINT ["PATH=/:$PATH ./module-runner"]
