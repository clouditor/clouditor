FROM node AS build

MAINTAINER Christian Banse <christian.banse@aisec.fraunhofer.de>

EXPOSE 80

WORKDIR /tmp

# this should hopefully trigger Docker to only update npm if dependencies have changed
ADD *.json ./
ADD *.lock ./
RUN yarn install --ignore-optional

# add the rest of the files
ADD . .

# set environment to production
ENV NODE_ENV production

# lint
RUN yarn lint

# build everything for production
RUN yarn run build --no-progress

FROM nginx:alpine

# copy to nginx
COPY --from=build /tmp/dist /usr/share/nginx/html/

ADD ./docker-entrypoint.sh /
ENTRYPOINT ["/bin/sh", "/docker-entrypoint.sh"]

EXPOSE 80
CMD ["nginx"]
