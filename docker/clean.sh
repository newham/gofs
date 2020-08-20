image="gofs"
version="2.0"

echo "->stop container"
docker stop $image

echo "->rm old docker"
docker rm $image
docker rmi $image:$version