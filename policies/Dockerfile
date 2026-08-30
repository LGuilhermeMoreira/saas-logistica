FROM openpolicyagent/opa:latest

WORKDIR /app

COPY policies/ /app/policies/

EXPOSE 8181

CMD ["run", "--server", "--addr=0.0.0.0:8181", "--config-file=/app/policies/config.yaml", "/app/policies/"]