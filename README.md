# call-tracker-backend
Backend part of the "Call Tracker"

## reuqest example
```
curl -v "ctrack.me/api/v1/phones/+79160000000" -X POST -F audio=@./c-dur.mp3
```
On success it should response with `201 Created`.
There is no check of the audio format sent.