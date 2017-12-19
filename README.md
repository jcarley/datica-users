
Build and package the binary into a docker container

    dc build package

Run the packaged up version of the binary

    dc run --rm -p 3000:3000 package


Ping the healthcheck

    curl http://localhost:3000/

