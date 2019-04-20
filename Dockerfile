FROM chrislu8468/alpine:1.0

MAINTAINER Chris <chrislu.name@gmail.com>

COPY snowflake /usr/bin/
RUN chmod +x /usr/bin/snowflake

ENTRYPOINT ["snowflake"]