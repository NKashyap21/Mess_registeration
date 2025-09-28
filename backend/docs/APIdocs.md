# Backend API Documentation


golang - gin
gorm
database - postgresql
API -> RestAPI


"message"
and
"error"
are on all endpoints


## students

POST /api/register?mess=0 ( 0 ,1 , 2 ,3 -> 4 messes)

request BODY
    userID

Response
    status code
        - 200 -> registration successful
        - 400 -> unsuccesful registration


- login -> oauth
jwt -> store user_id


GET /api/user
request BODY
    "jsontoken"
    or oauth middle ware
    verify this

    status code
        - 200 OK
        - 401 Unauthorized
    message
        - "succes"
        - "userType"
        - "name"
        - "rollno"
        - "veg/non-veg"
        - "mess"

### swapping request

ONLY FOR A TO B OR B TO A
no inner swapping


1. user posts for request
    -  friends
    -  public
2. show notificatoins
    - when public requests are auto swapped
    - when friend requests are approved
3. able to see all requests rightnow

POST /api/swap-request
    - authorized via jwt

    request BODY
        - name
        - email
        - type (friend, public)
        - user_id
    response
        status code
            - 201 created
            - 401 unauthorized
            - 400 bad request
        message:

DELETE /api/swap-request
    - authorized via jwt

    request BODY
        - name
        - email
    response
        status code
            - 201 created
            - 401 unauthorized
            - 400 bad request
        message:


GET /api/get-swaps
    - authorized via jwt

    request BODY
        empty
    response
        - array of swaps
        both sides
        "AtoB":[]
        "BtoA":[]


### scanning for mess people

api - key for mess people
so that no one messes up with this endpoint

GET /api/scanning?rollno=es23btech1028

request BODY
    empty
Response
    status code
        - 200 -> registration successful
        - 400 -> unsuccesful registration
    message
        - "succes"
        - "userType"
        - "name"
        - "rollno"
        - "veg/non-veg"

- things to add
    special dinner
    day passes

---

## hostel office

- sql viewer

GET all users -> google contacts kinda

GET /api/user?rollno=es23btech1028
    request BODY
    somw jwt auth

    response
        status code
            - 200 -> registration successful
            - 400 -> unsuccesful registration
        message
            - "succes"
            - "userType"
            - "name"
            - "rollno"
            - "veg/non-veg"
            - "mess"
            ... what ever we store

may be can track users swaps in breakfast , lunch , snacsk , dinner

POST /api/update
GET /api/user?rollno=es23btech1028
    somw jwt auth

    request BODY
     ... any info here is what will reflect in db
    response
        status code
            - 200 -> registration successful
            - 400 -> unsuccesful registration
        message
            - "succes"
            - "userType"
            - "name"
            - "rollno"
            - "veg/non-veg"
            - "mess"
            ... what ever we store

- will talk to hostel office about this
- mess wise analytics
- student wise
    - btech
    - mtech
    - phd


## to add

reset endpoint to reset all data for users when registration starte
download .csv


## after this apis

- table design will be added soon...
