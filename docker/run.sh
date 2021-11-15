#docker run -v /vagrant_data/http-files:/gohttpserver/http-files/ --name go-http-fileserver -d gohttpserver:latest 
image="gofs"
version="1.0"
outport=8087
inport=8087
savedfilepath=$(dirname "$PWD")/files

echo "rm old docker"
docker stop $image
docker rm $image
docker rmi $image:$version

#echo "build"
#cd ../
#go build --ldflags "-extldflags -static"
#go build --ldflags '-extldflags "-static -lstdc++ -lpthread"'
#go build
#cd docker

mkdir copy
mkdir copy/files
cp -r ../conf ../public ../view ../LICENSE ../gofs ./copy;

echo "build docker"
docker build -t $image:$version .
echo "run docker"
docker run --name $image -p $outport:$inport --restart=always -v $savedfilepath:/gofs/files -d $image:$version

rm -rdf copy

