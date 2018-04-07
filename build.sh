docker build -t call-tracker-backend .
# docker tag 5b3ff735943a masim05/call-tracker-backend
# docker push masim05/call-tracker-backend
docker run -p 80:80 -d --name call-tracker-backend masim05/call-tracker-backend