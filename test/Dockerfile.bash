FROM bash:latest

RUN apk update && apk add --no-cache curl imagemagick jpeg

COPY --chmod=755 ./scripts/integration_test.sh /scripts/integration_test.sh
COPY ./images/modified ./modified_images

CMD ["/scripts/integration_test.sh"]