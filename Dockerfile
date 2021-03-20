ARG TAG=latest
# always fetch the latest ui for now. later we might pin specific versions
FROM clouditor/ui:latest AS dashboard

FROM clouditor/engine:${TAG}

COPY --from=dashboard /usr/share/nginx/html /usr/local/clouditor/html/

CMD ["--db-in-memory"]
