MODULE=csgeo
GOOS=linux go build -o $MODULE .
docker build -t $MODULE .
docker tag $MODULE localhost:5000/$MODULE
docker push localhost:5000/$MODULE
docker rmi localhost:5000/$MODULE
docker pull localhost:5000/$MODULE

