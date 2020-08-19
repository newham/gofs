#docker run -v /vagrant_data/http-files:/gohttpserver/http-files/ --name go-http-fileserver -d gohttpserver:latest 
image="gofs"
version="2.0"
outport=8087
inport=8087
savedfilepath=$HOME/http-files

echo "build"
cd ../
#go build --ldflags "-extldflags -static"
#go build --ldflags '-extldflags "-static -lstdc++ -lpthread"'
go build
cd docker

mkdir copy
mkdir copy/files
cp -r ../conf ../public ../gofs ./copy;

echo "stop container"
docker stop $image

echo "rm old docker"
docker rm $image
docker rmi $image:$version

echo "build docker"
docker build -t $image:$version .
echo "run docker"
docker run --name $image -p $outport:$inport -v $savedfilepath:/gofs/files -d $image:$version

rm -rdf copy

