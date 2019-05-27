ARG TAG=latest
FROM clouditor/ui:${TAG} AS dashboard

FROM clouditor/engine:${TAG}

COPY --from=dashboard /usr/share/nginx/html /usr/local/clouditor/html/

CMD ["--db-in-memory"]
