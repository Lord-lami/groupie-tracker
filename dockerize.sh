#!/bin/bash

# Build the image
imagename="sheeplami/groupie-tracker:1.0.1"
echo "BUILDING THE $imagename IMAGE"
docker image build -f Dockerfile --label "env=dev" -t $imagename .

containername="groupie-tracker"

# Remove the previous container if it exists
docker rm -f $containername || echo "$containername does not exist. Making a new one"

# Run the container

echo "STARTING THE $containername CONTAINER OFF $imagename IMAGE"
docker container run -p 8080:8080 --label "env=dev" --detach --name $containername $imagename