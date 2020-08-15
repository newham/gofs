#docker run -v /vagrant_data/http-files:/gohttpserver/http-files/ --name go-http-fileserver -d gohttpserver:latest 
echo "save image"
docker save $image:$version>$image.tar

rm -rdf copy

