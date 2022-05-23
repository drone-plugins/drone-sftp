FROM alpine:3.16

RUN apk update && \
  apk add \
    ca-certificates \
    git \
    nodejs && \
  rm -rf \
    /var/cache/apk/*

WORKDIR /node
ADD package.json /node/
ADD index.js /node/
RUN npm install --production

ENTRYPOINT ["node", "index.js"]
