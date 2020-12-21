# ChatAPI 

Run: `docker-compose up` (first time, add --build flag)


API port:       8080


## API endpoints


`/auth/reg`
~~~~
Description:
    Registers new user and returns authorization token with some user data

Protocol:
    http

Method:
    POST

Request data:
    email: string
    username: string
    password: string

Response data:
    token: string
    username: string
    id: int
~~~~

`/auth/login`
~~~~
Description:
    User authorization, returns token and some user data

Protocol:
    http

Method:
    POST

Request data:
    email: string
    password: string

Response data:
    token: string
    username: string
    id: int
~~~~

`/users/me`
~~~~
Description:
    Returns data of current user (requires token)

Protocol:
    http

Method:
    GET

Permissions:
    Only authorized users

Response data:
    id: int
    email: string
    username: string
    password (hashed): string
    is_online: bool
~~~~

`/users/`
~~~~
Description:
    Returns all users from database

Protocol:
    http

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
    is_online: bool
~~~~

`/chats/create`
~~~~
Description:
    Creates new chat (requires token)

Protocol:
    http

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
Description:
    Returns all messages from given chat 

Protocol:
    http

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

`/chats/ws?token=`
~~~~
Description:
    Websocket connection for sending and receiving messages (requires token).
    Current user goes online when connection opens, and goes offline when closes.

Protocol:
    Websocket

Permissions:
    Only authorized users

Request data:
    text: string
    chat_id: int

Response data:
    text: string
    username: string
    created_at: time
~~~~