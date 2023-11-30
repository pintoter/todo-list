# todo-list

## Examples

Create note:
$ curl -X 'POST' \
  'http://127.0.0.1:8080/api/v1/note' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "date": "2023-11-30",
  "description": "Description",
  "status": "done",
  "title": "Note "
}'

