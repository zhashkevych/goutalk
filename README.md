# goutalk
Simple chat server written in Go (demo project for "Backend Challenge Hackathon") 

## Requirements
- docker
- docker-compose

## Build and Run

To run the project locally all you have to do is run:

```docker-compose up server```

## Key Concept and Endpoints description

###POST /login

Requires no auth, creates new user if not exists, or sings in existent user by generating JWT Token 

##### Example Input: 
```
{
	"user_name": "conorMcGregor",
	"password": "illbeatyoass"
} 
```

##### Example Response:
```
{
    "id": "5d091716db8212e1a2efb33a",
    "user_name": "conorMcGregor",
    "credentials": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjA5NjMyMjUsImlhdCI6MTU2MDg3NjgyNSwidXNlcl9pZCI6IjVkMDkxNzE2ZGI4MjEyZTFhMmVmYjMzYSIsInVzZXJuYW1lIjoibmV3dXNlciIsInBhc3N3b3JkIjoiXHVmZmZkXHVmZmZkJ1x1MDAxMlx1ZmZmZFx1ZmZmZEBcdWZmZmRcdTAwMTJRXHVmZmZkXHVmZmZkYdWt3qVv25oifQ.eruNfiyxFSm3H1s4uleY9Cuxw-D6qbjukjwI1kTwn70"
}
```

##### *Note: All other endpoints requires Bearer token authentication

###GET /users

Returns a list of existing users

##### Example Response:
```
{
    "users": [
        {
            "user_id": "5d061d45c73985d0b598b0a5",
            "user_name": "nagibator"
        },
        {
            "user_id": "5d091e4ddb8212e1a2efb33c",
            "user_name": "conorMcGregor"
        }
    ]
}
```

###GET /users/<USER_ID>

Returns info for specific user

##### Example Response:
```
{
    "user_id": "5d091e4ddb8212e1a2efb33c",
    "user_name": "conorMcGregor"
}
```

###GET /rooms

Returns a list of rooms available

##### Example Response:
```
[
    {
        "room_id": "5d07348c594a0972c8010faf",
        "room_name": "memes",
        "creator_id": "5d062c9120f31443e4cdd449",
        "created_at": "2019-06-17T06:34:52.504Z",
        "members": [
            {
                "user_id": "5d091e4ddb8212e1a2efb33c",
                "user_name": "conorMcGregor"
            }
        ]
    },
    {
        "room_id": "5d07388c32749f0d3885cb04",
        "room_name": "morememes",
        "creator_id": "5d062c9120f31443e4cdd449",
        "created_at": "2019-06-17T06:51:56.005Z",
        "members": []
    }
]
```

###GET /rooms/<ROOM_ID>

Returns info for a given room

##### Example Response:
```
 {
    "room_id": "5d07348c594a0972c8010faf",
    "room_name": "memes",
    "creator_id": "5d062c9120f31443e4cdd449",
    "created_at": "2019-06-17T06:34:52.504Z",
    "members": [
        {
            "user_id": "5d091e4ddb8212e1a2efb33c",
            "user_name": "conorMcGregor"
        }
    ]
}
```

###POST /rooms

Adds a new chat room

##### Example Input: 
```
{
	"name": "bestoftinder"
}
```

##### Example Response:
```
{
    "room_id": "5d092005db8212e1a2efb33d",
    "room_name": "bestoftinder",
    "creator_id": "5d091716db8212e1a2efb33a",
    "created_at": "2019-06-18T17:31:49.428210839Z",
    "members": []
}
```

###DELETE /rooms/<ROOM_ID>

Removes specific chat room (by room creator)

##### Example Response:
```
{
    "message": "room removed successfully"
}
```

###POST /rooms/<ROOM_ID>/join

Adds current user to a given chatroom

##### Example Response:
```
{
    “result”:  “user successfully joined the room”
}
```


###POST /rooms/<ROOM_ID>/leave

Removes current user from a given chatroom

##### Example Response:
```
{
    “result”:  “user successfully left the room”
}
```

###POST /message

Sends message to all connected via Websockets clients
##### *Websoket connection is established on ws://hostname:1030/

##### Example Input: 
```
{
	"room_id": "5d07348c594a0972c8010faf",
	"message": "I’m going to the stars and then past them."
}
```

##### Example Response:
```
{
    "result": "message successfully sent" 
}
```

##### *Clients connected via Websocket recieves payload like this:
```
{
   "user_id": "5d092005db8212e1a2efb33d", 
   "room_id": "5d07348c594a0972c8010faf",
   "message": "I’m going to the stars and then past them."
}
```