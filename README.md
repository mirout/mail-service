# Mail Service

This is a simple mail service that can be used to send emails.

## Running the service

To run the service, you need to have a docker compose:
```bash
docker compose up
```

You should specify the following environment variables:
- `POSTGRES_PASSWORD` - password for the postgres database
- `REDIS_PASSWORD` - password for the redis database
- `MAIL_USERNAME` - username for the mail service
- `MAIL_PASSWORD` - password for the mail service

Also you can edit the `docker-compose.yml` file to change the following args:
- `--smtp-host` - host for the smtp server
- `--smtp-port` - port for the smtp server
- `--db-host` - host for the postgres database
- `--db-port` - port for the postgres database default is 5432
- `--redis-host` - host for the redis database
- `--redis-port` - port for the redis database default is 6379
- `--mail-host` - host for the mail service with protocol (for example `http://localhost:8080`)

## Usage

### Handlers

All requests should be sent to the `/api/v1` endpoint except the `/img` endpoint.

All time fields should be in the RFC3339 format.

#### `/users` endpoint
To register a new user, you need to send a POST request to `/api/v1/users` with the following body:
```json5
{
    "email": "email@example.com",
    "first_name" : "First Name",
    "last_name" : "Last Name"
}
```
It will return a response with id of the user:
```
7e2c026b-32b6-4957-94a3-b08b0242b213
```

To get a user, you need to send a GET request to `/api/v1/users` with one of the following query params:
- `id` - id of the user
- `email` - email of the user

It will return a response with the user:
```json5
{
    "id": "7e2c026b-32b6-4957-94a3-b08b0242b213",
    "email": "email@example.com",
    "first_name" : "First Name",
    "last_name" : "Last Name",
    "created_at": "2021-09-05T12:00:00Z"
}
```

#### `/groups` endpoint

To create a new group, you need to send a POST request to `/api/v1/groups` with the following body:
```json5
{
    "name": "Group Name"
}
```

It will return a response with id of the group:
```
7e2c026b-32b6-4957-94a3-b08b0242b213
```

To get a group, you need to send a GET request to `/api/v1/groups` with id of the group in the query params. It will return a response with the group:
```json5
{
    "id": "7e2c026b-32b6-4957-94a3-b08b0242b213",
    "name": "Group Name",
    "created_at": "2021-09-05T12:00:00Z"
}
```

To add a user to a group, you need to send a POST request to `/api/v1/groups/{group_id}/add/{user_id}` with the empty body.
To remove a user from a group, you need to send a POST request to `/{group_id}/remove/{user_id}` with the empty body.

#### `/mails` endpoint

To send a mail to user, you need to send a POST request to `/api/v1/mails/send/to/user/{user_id}` with the following body:
```json5
{
    "subject": "Subject",
    "body": "Body",
    "send_at": "2021-09-05T12:00:00Z" // optional field to send mail at a specific time
}
```

To send a mail to group, you need to send a POST request to `/api/v1/mails/send/to/group/{group_id}` with the following body:
```json5
{
    "subject": "Subject",
    "body": "Body",
    "send_at": "2021-09-05T12:00:00Z" // optional field to send mail at a specific time
}
```

To get a mail, you need to send a GET request to `/api/v1/mails/{mail_id}`. It will return a response with the mail:
```json5
{
    "id": "7e2c026b-32b6-4957-94a3-b08b0242b213",
    "subject": "Subject",
    "body": "Body",
    "sent_at": "2021-09-05T12:00:00Z",
    "created_at": "2021-09-05T12:00:00Z"
}
```

To get all mails which was sent to user, you need to send a GET request to `/api/v1/mails/to/user/{user_id}`. It will return a response with the list of mails:
```json5
[
    {
        "id": "7e2c026b-32b6-4957-94a3-b08b0242b213",
        "subject": "Subject",
        "body": "Body",
        "sent_at": "2021-09-05T12:00:00Z",
        "created_at": "2021-09-05T12:00:00Z"
    }
]
```

#### `/img` endpoint

This endpoint is used to get an 1x1 image to track if the email was opened. To get the image, you need to send a GET request to `/img/{mail_id}`.

### Templates

You can use the following templates in the body of the mail:
- `{{.FirstName}}` - first name of the user
- `{{.LastName}}` - last name of the user
- `{{.Body}}` - body of the mail