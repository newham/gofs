#docker run -v /vagrant_data/http-files:/gohttpserver/http-files/ --name go-http-fileserver -d gohttpserver:latest 
image="gofs"
version="2.0"
outport=8087
inport=8087
# 默认根目录下files
savedfilepath=$(dirname "$PWD")/files
echo "\n->savedfilepath:$savedfilepath\n"

echo "\n->build go bin\n"
cd ../
#go build --ldflags "-extldflags -static"
#go build --ldflags '-extldflags "-static -lstdc++ -lpthread"'
go build
if [ $? -ne 0 ]; then
    echo "build go bin failed, exit"
    exit 1

echo "\n->start to build docker\n"

cd docker

mkdir copy
mkdir copy/files
cp -r ../conf ../public ../gofs ./copy;

echo "\n->stop container\n"
docker stop $image

echo "\n->rm old docker\n"
docker rm $image
docker rmi $image:$version

echo "\n->build docker\n"
docker build -t $image:$version .
echo "run docker"
docker run --name $image -p $outport:$inport -v $savedfilepath:/gofs/files -d $image:$version

rm -rdf copy

