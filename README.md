```
docker run --name db -e POSTGRES_USER=? -e POSTGRES_PASSWORD=? -p 4444:5432 -d postgres:alpine
docker exec -it db
psql -U USER?
CREATE DATABASE db;
```

```
migrate -database 'postgres://user?:pass?@0.0.0.0:4444/db?sslmode=disable' -source file://? -verbose up
migrate -database 'postgres://user?:pass?@0.0.0.0:4444/db?sslmode=disable' -source file://? -verbose down
```

```
docker build -t asia.gcr.io/project-id/container-bucket .
docker push asia.gcr.io/project-id/container-bucket
gcloud auth configure-docker
```

```ENV
APP_HOST=
APP_PORT=
APP_NAME=
APP_VERSION=
APP_BODY_LIMIT=
APP_READ_TIMEOUT=
APP_WRTIE_TIMEOUT=
APP_FILE_LIMIT=
APP_GCP_BUCKET=

JWT_ADMIN_KEY=
JWT_SECRET_KEY=
JWT_API_KEY=
JWT_ACCESS_EXPIRES=
JWT_REFRESH_EXPIRES=

DB_HOST=
DB_PORT=
DB_PROTOCOL=
DB_USERNAME=
DB_PASSWORD=
DB_DATABASE=
DB_SSL_MODE=
DB_MAX_CONNECTIONS=
```
