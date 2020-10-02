FROM alpine
EXPOSE 53589

RUN apk update && apk add taskd
USER taskd

WORKDIR /var/lib/taskd
RUN mkdir data && mkdir certs
COPY docker-images/config /var/lib/taskd/data/config

ENTRYPOINT taskd 
CMD server --data /var/lib/taskd/data
