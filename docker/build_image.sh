#docker run -v /vagrant_data/http-files:/gohttpserver/http-files/ --name go-http-fileserver -d gohttpserver:latest 
image="gofs"
version="1.0"
outport=8087
inport=8087
#savedfilepath=/home/local/http-files

echo "build"
cd ../
#go build --ldflags "-extldflags -static"
#go build --ldflags '-extldflags "-static -lstdc++ -lpthread"'
go build
cd docker

mkdir copy
mkdir copy/files
cp -r ../conf ../public ../view ../LICENSE ../gofs ./copy;

echo "rm old docker"
docker rm $image
docker rmi $image:$version

echo "build docker"
docker build -t $image:$version .

echo "save image"
docker save $image:$version>$image.tar

rm -rdf copy

