# ChatAPI 

Run: `docker-compose up` (first time, add --build flag)


API port:       8080

Websocket port: 8000

# Websocket endpoint
~~~~
Request data:
    text: string
    chat_id: int
    owner_id: int

Response data:
    id: int
    text: string
    chat_id: int
    owner_id: int
~~~~


## API endpoints


`/auth/reg`
~~~~
Method:
    POST

Request data:
    id: int
    email: string
    username: string
    password: string

Response data:
    token: string
~~~~

`/auth/login`
~~~~
Method:
    POST

Request data:
    id: int
    email: string
    password: string

Response data:
    token: string
~~~~

`/users/me`
~~~~
Method:
    GET

Permissions:
    Only authorized users

Response data:
    id: int
    email: string
    username: string
    password (hashed): string
~~~~

`/users/`
~~~~
Method:
    GET

Permissions:
    Only authorized users

Response data:
Array of:
    id: int
    email: string
    username: string
    password (hashed): string
~~~~

`/chats/create`
~~~~
Method:
    POST

Permissions:
    Only authorized users

Request data:
    name: string

Response data:
    id: int
    name: string
~~~~

`/chats/messages?id=`
~~~~
Method:
    GET

Permissions:
    Only authorized users

Request query data:
    id: int

Response data:
Array of:
    id: int
    text: string
    owner_id: int
    chat_id: int
    created_at: time
~~~~