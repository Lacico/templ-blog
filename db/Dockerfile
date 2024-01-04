FROM alpine:20231219

RUN apk update && apk add sqlite --no-interactive
RUN mkdir /db/
WORKDIR /db/ 
RUN sqlite3 blogDB.db
CMD /bin/sh


