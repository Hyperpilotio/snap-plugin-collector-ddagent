FROM alpine:3.4

EXPOSE 8182
EXPOSE 8000

COPY snap-ddagent /usr/local/bin/

ENTRYPOINT ["snap-ddagent", "--stand-alone", "--stand-alone-port", "8182"]

CMD [""]

