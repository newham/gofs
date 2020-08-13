var fileCount = 0;
var uploader

function set_upload(path, f) {
    uploader = new qq.FineUploader({
        autoUpload: true,
        debug: false,
        element: document.getElementById('fine-uploader'),
        request: {
            endpoint: '/upload',
            inputName: 'file',
            params: {
                path: path
            }
        },
        validation: {
            itemLimit: 20
        },
        deleteFile: {
            enabled: false,
            endpoint: ''
        },
        retry: {
            enableAuto: false
        },
        callbacks: {
            onAllComplete: function(succeeded, failed) {
                if (failed.length == 0) {
                    console.log("Upload Success");
                    f()
                } else {
                    console.log("Upload Failed");
                }
                // window.location.href = "/folder?name={{$.Folder.Path}}";
            },
            onSubmitted: function(id, name) {
                // $('#upload-info').text(++fileCount);
                // $('#upload-info').show();
            },
            onCancel: function(id, name) {
                // $('#upload-info').text(--fileCount);
                // if (fileCount == 0) {
                //     $('#upload-info').hide();
                // }

            }
        }
    });
}

function doUpload() {
    uploader.uploadStoredFiles();
}