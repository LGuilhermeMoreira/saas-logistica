FROM openpolicyagent/opa:latest

WORKDIR /app

COPY policies/ /app/policies/

COPY config.yaml /app/config.yaml

EXPOSE 8181

CMD ["run", "--server", "--config-file=/app/config.yaml", "/app/policies/"]