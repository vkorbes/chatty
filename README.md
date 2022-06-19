(This was a job interview exercise a billion years ago.)

This API will be used to control users exchanging messages between each other.

The command takes two optional arguments:

- `-p` Is the the port number. Default is 8000.
- `-m` Is the MongoDB URL. Default is whatever's in your secret.txt file on the same folder as the executable.

Here are some things you can do with this app:

- POST request to `[URL]/users` containing `{"name": "User Name","username": "username"}` adds that entry to the database.

Example output:
```
{
    "id": "5a92fe5a7d9b532f98e8bba1",
    "budget": 10,
    "name": "User Name",
    "username": "username",
    "createdAt": "2018-02-25T18:20:10.805Z",
    "updatedAt": "2018-02-25T18:20:10.805Z"
}
```

- GET request to `[URL]/users/[User ID]` gets a user from the database. For example, after the request above has been processed, a request to `[URL]/users/5a92fe5a7d9b532f98e8bba1` would yield the same output.

- POST request to `[URL]/messages` containing `{"from": "orange","to": "banana","body": "This is a test message."}` adds that message to the database.

Example output:
```
{
    "id": "5a93000c7d9b532f98e8bba2",
    "from": "orange",
    "to": "banana",
    "body": "This is a test message.",
    "sentAt": "2018-02-25T18:27:24.885Z"
}
```

- GET request to `[URL]/message/[Message ID]` gets a message from the database. For example, after the request above has been processed, a request to `[URL]/messages/5a93000c7d9b532f98e8bba2` would yield the same output.

- GET request to `[URL]/messages?to=username` gets all messages addressed to that username.

Example output:
```
{
    "messages": [
        {
            "id": "5a8d76807d9b537448d19b30",
            "from": "orange",
            "to": "banana",
            "body": "This is a test message.",
            "sentAt": "2018-02-21T13:39:12.767Z"
        },
        {
            "id": "5a93000c7d9b532f98e8bba2",
            "from": "apple",
            "to": "banana",
            "body": "This is another test message.",
            "sentAt": "2018-02-25T18:27:24.885Z"
        }
    ]
}
```

- GET request to `[URL]/listusers` lists all users. This is not on spec, it's there as a development aid.

Example output:

```
[
    {
        "id": "5a8d75057d9b53706595116a",
        "budget": 9,
        "name": "Orange",
        "username": "orange",
        "createdAt": "2018-02-21T13:32:53.509Z",
        "updatedAt": "2018-02-21T13:39:13.159Z"
    },
    {
        "id": "5a8d750d7d9b53706595116b",
        "budget": 8,
        "name": "Banana",
        "username": "banana",
        "createdAt": "2018-02-21T13:33:01.239Z",
        "updatedAt": "2018-02-21T13:38:52.757Z"
    }
]
```

- GET request to `[URL]/listmsg` lists all messages. This is not on spec, it's there as a development aid.

Example output:

```
[
    {
        "id": "5a8d766c7d9b537448d19b2f",
        "from": "banana",
        "to": "orange",
        "body": "Message.",
        "sentAt": "2018-02-21T13:38:52.358Z"
    },
    {
        "id": "5a93000c7d9b532f98e8bba2",
        "from": "orange",
        "to": "banana",
        "body": "This is a test message.",
        "sentAt": "2018-02-25T18:27:24.885Z"
    }
]
```
