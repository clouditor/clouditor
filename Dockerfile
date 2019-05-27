ARG TAG=latest
FROM clouditor/clouditor-dashboard:${TAG} AS dashboard

FROM clouditor/clouditor-engine:${TAG}

COPY --from=dashboard /usr/share/nginx/html /usr/local/clouditor/html/

CMD ["--db-in-memory"]
