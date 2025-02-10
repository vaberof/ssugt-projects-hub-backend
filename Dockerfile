FROM golang:1.22

EXPOSE 8080

RUN mkdir -p /opt/app

COPY cmd/ssugt-projects-hub/build /opt/app

WORKDIR /opt/app

CMD [ "opt/app/main", "-config.files", "container.yaml", "-env.vars.file", "application.env" ]
