# one time
# create volume
docker volume create call_tracker_pgdata
# everytime
docker build -t call-tracker-db .
docker run -p 5432:5432 -d -v call_tracker_pgdata:/var/lib/postgresql/data --name call_tracker_db call-tracker-db
