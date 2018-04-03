# call-tracker-backend
Backend part of the "Call Tracker"

## Request examples
### Save call record
```
curl -v "ctrack.me/api/v1/phones/+79160000000" -X POST -F audio=@./c-dur.mp3
```
On success it should response with `201 Created`.
There is no check of the audio format sent.

### Get call metrics
```
curl ctrack.me/api/v1/calls/10
> {"stt":{"text":"Пусть просто такую информацию легче записать голосом чем писать текстом"},"duration":{"duration":4.56}}

curl ctrack.me/api/v1/calls/11
> {"stt":{"text":"Пусть просто такую информацию легче записать голосом чем писать текстом"}}

curl ctrack.me/api/v1/calls/12
> {"duration":{"duration":4.56}}
```
Call 10 has both metrics calculated, call 11 has only stt and call 12
has only duration metric.
On success 200 OK will be returned.
The response contains pairs: `metric name`: `result of the correspondent
endpoint`. More information about available metrics are [here](https://github.com/vladimir-chernykh/call-tracker-processing/tree/master/code/endpoints).