FROM liquibase/liquibase:5.0
RUN lpm add postgresql --global
COPY ./changelog /liquibase/changelog
COPY ./liquibase.properties /liquibase/liquibase.properties

WORKDIR /liquibase

ENTRYPOINT [ "liquibase" ]
CMD ["update"]