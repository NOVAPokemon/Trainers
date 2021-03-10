FROM brunoanjos/nova-server-base:latest

ENV executable="executable"
COPY $executable .

EXPOSE 8009

CMD ["sh", "-c", "./$executable -a -l"]