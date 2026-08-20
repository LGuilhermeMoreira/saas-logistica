# FROM openpolicyagent/opa:latest

# WORKDIR /app

# COPY policies/ /app/policies/

# EXPOSE 8181

# CMD ["run", "--server", "--config-file=/app/config.yaml", "/app/policies/"]

FROM openpolicyagent/opa:latest

WORKDIR /app

COPY policies/ /app/policies/

EXPOSE 8181

CMD ["run", "--server", "--config-file=/app/policies/config.yaml", "/app/policies/"]