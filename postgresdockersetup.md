# PostgreSQL Docker Setup
1. Pull the postrges image - docker pull postgres
2. Create a volume to persist data - mkdir -p <filepath>
3. Run docker run --rm --name postgres-instance -e POSTGRES_USER=go -e POSTGRES_PASSWORD=gopassword -d -p 5432:5432 -v /Users/rhysjohns/Code/Docker/Volumes/postgres:/var/lib/postgresql/data postgres
4. Connect via psql postgresql://<username>:<password>@<host>:<port>/<database>