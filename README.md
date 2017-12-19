## Prerequisites

  Tested and developed using the following specs.  You'll need a docker
  host to run the example application.

    * Mac OSX 10.12.6
    * Docker for Mac 17.09.1
    * curl
    * Go 1.9.2 (This isn't entirely necessary unless running the app locally)

## Getting started

  To start up the application run the following commands

    make database
    make package
    make run

  Test that everything is up and running with a request to the healthcheck.  You should
  see the Hello World greeting

    curl http://localhost:3000/


  It's helpful if you pipe the results of the curl requests through jq.

    brew install jq

    curl -s http://localhost:3000/ | jq .

  Control + C quits the application gracefully

## API Examples


### Signup a user

    curl -s -H "Content-Type: application/json" -X POST -d '{"email":"jeff.carley@example.com","password":"secret", "name":"Jeff Carley"}' http://localhost:3000/user

### Signin a user

    curl -s -H "Content-Type: application/json" -X POST -d '{"username":"jeff.carley@example.com","password":"secret"}' http://localhost:3000/auth

### Get a given user's info

    curl -s -H "Content-Type: application/json" -H "Authorization: bearer {signature}" -XGET -d '' http://localhost:3000/user/jeff.carley@example.com

  _Example_:

    curl -s \
      -XGET \
      -H "Content-Type: application/json" \
      -H "Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJmMTNmOGEwNi01MWY5LTRiYTctYjAzNC02M2YwNGZjNmU0NzMiLCJ1c2VybmFtZSI6ImplZmYuY2FybGV5QGV4YW1wbGUuY29tIiwiaWF0IjoiMjAxNy0xMi0xOVQxNjowODoxNVoifQ.ASdFIBKHTGTPI1Tq9oguCcSJezXeMws-LLzFr6VFMAU" \
      http://localhost:3000/user/jeff.carley@example.com

### Update a given user's info

    curl -s -H "Content-Type: application/json" -H "Authorization: bearer {signature}" -XPUT -d '' http://localhost:3000/user/jeff.carley@example.com

  _Example_:

    curl -s \
      -XPUT \
      -H "Content-Type: application/json" \
      -H "Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJmMTNmOGEwNi01MWY5LTRiYTctYjAzNC02M2YwNGZjNmU0NzMiLCJ1c2VybmFtZSI6ImplZmYuY2FybGV5QGV4YW1wbGUuY29tIiwiaWF0IjoiMjAxNy0xMi0xOVQxNjowODoxNVoifQ.ASdFIBKHTGTPI1Tq9oguCcSJezXeMws-LLzFr6VFMAU" \
      -d '{"name":"John Doe"}' \
      http://localhost:3000/user/jeff.carley@example.com

### Delete a given user's info

    curl -s -H "Content-Type: application/json" -H "Authorization: bearer {signature}" -XDELETE http://localhost:3000/user/jeff.carley@example.com

  _Example_:

    curl -s \
      -XDELETE \
      -H "Content-Type: application/json" \
      -H "Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJmMTNmOGEwNi01MWY5LTRiYTctYjAzNC02M2YwNGZjNmU0NzMiLCJ1c2VybmFtZSI6ImplZmYuY2FybGV5QGV4YW1wbGUuY29tIiwiaWF0IjoiMjAxNy0xMi0xOVQxNjowODoxNVoifQ.ASdFIBKHTGTPI1Tq9oguCcSJezXeMws-LLzFr6VFMAU" \
      http://localhost:3000/user/jeff.carley@example.com



## Common development tasks

Build and package the binary into a docker container

    dc build package

Run the packaged up version of the binary

    dc run --rm -p 3000:3000 package

Ping the healthcheck

    curl http://localhost:3000/

Remove dangling images

    docker rmi $(docker images --quiet --filter "dangling=true")


