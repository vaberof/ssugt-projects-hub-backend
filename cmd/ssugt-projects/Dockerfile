FROM golang:latest

EXPOSE 8080

WORKDIR /opt/app

COPY cmd/ssugt-projects/build /opt/app

ENTRYPOINT [ "./main", "-config.files", "application.yaml", "-env.vars.file", "application.env" ]
