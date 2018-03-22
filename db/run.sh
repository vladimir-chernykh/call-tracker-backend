# one time
# create volume
docker volume create call_tracker_pgdata
# everytime
docker build -t call-tracker-db .
docker run -p 172.17.0.1:5432:5432 -d -v call_tracker_pgdata:/var/lib/postgresql/data --name call-tracker-db masim05/call-tracker-db
