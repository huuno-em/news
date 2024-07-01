FROM postgres:10.0-alpine

ENV POSTGRES_USER test_user
ENV POSTGRES_PASSWORD qwerty123
ENV POSTGRES_DB testdb

# Копирование файла schema.sql внутрь контейнера
COPY pcg/database/schema.sql /docker-entrypoint-initdb.d/

# Установка правильных разрешений на файл
RUN chmod 755 /docker-entrypoint-initdb.d/schema.sql

# Открытие порта для доступа к PostgreSQL
EXPOSE 5432

# Запуск PostgreSQL при запуске контейнера
CMD ["postgres"]