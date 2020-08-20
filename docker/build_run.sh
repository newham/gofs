#docker run -v /vagrant_data/http-files:/gohttpserver/http-files/ --name go-http-fileserver -d gohttpserver:latest 
image="gofs"
version="2.0"
outport=8087
inport=8087
# 默认根目录下files
savedfilepath=$(dirname "$PWD")/files
echo "->savedfilepath:$savedfilepath"

echo "->build go bin"
cd ../
#go build --ldflags "-extldflags -static"
#go build --ldflags '-extldflags "-static -lstdc++ -lpthread"'
# 升级依赖包
go get -u github.com/newham/hamgo 
go build
if [ $? -ne 0 ]; then
    echo "build go bin failed, exit"
    exit 1
fi

echo "->start to build docker"

cd docker

mkdir copy
mkdir copy/files
cp -r ../conf ../public ../gofs ./copy;

echo "->stop container"
docker stop $image

echo "->rm old docker"
docker rm $image
docker rmi $image:$version

echo "->build docker"
docker build -t $image:$version .
echo "run docker"
docker run --name $image -p $outport:$inport -v $savedfilepath:/gofs/files -d --restart=always $image:$version

rm -rdf copy