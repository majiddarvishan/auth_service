FROM router_core:1.12.7

WORKDIR /app

COPY auth_service main

RUN rm -rf docs
COPY  docs docs

EXPOSE 8080
EXPOSE 8443

CMD ["./main"]
