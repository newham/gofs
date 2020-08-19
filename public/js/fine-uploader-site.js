var fileCount = 0;
var uploader

function set_upload() {
    uploader = new qq.FineUploader({
        autoUpload: true,
        debug: false,
        element: document.getElementById('fine-uploader'),
        request: {
            endpoint: '/upload',
            inputName: 'file',
            // params: {
            //     dir: get_current(),
            // }
        },
        validation: {
            itemLimit: 256
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
                    reload_page()
                } else {
                    console.log("Upload Failed");
                }

                fileCount = 0 // 上传数量清零
                $('#task-box').hide()

                // clean_done_list() // 注释掉则是手动清空
                // window.location.href = "/folder?name={{$.Folder.Path}}";
                // show_box_bg(true)
            },
            onComplete: function(id, name) {
                --fileCount
                $('#task-count').html(fileCount);
            },
            onSubmitted: function(id, name) {
                ++fileCount
                $('#task-count').html(fileCount);
                // $('#upload-info').text(++fileCount);
                // $('#upload-info').show();
                // show_box_bg(false)
            },
            onCancel: function(id, name) {
                --fileCount
                $('#task-count').html(fileCount);
                // $('#upload-info').text(--fileCount);
                // if (fileCount == 0) {
                //     $('#upload-info').hide();
                // }
            }
        }
    });
}

function show_box_bg(show) {
    if (show) {
        $('#qq-uploader').css("background-image", `url("/public/img/svg/icon-collect-info.svg") 100% 100% no-repeat;`)
    } else {
        $('#qq-uploader').css("background", "white")
    }
}

function doUpload() {
    uploader.uploadStoredFiles();
}

function clean_done_list() {
    if (fileCount > 0) {
        alert('还有未完成的任务！')
        return
    }
    uploader.clearStoredFiles()
    show_box_bg(true)
}

function set_upload_path(path) {
    // 如果有上传的任务，禁止修改path！
    if (fileCount > 0) {
        return
    }
    // console.log('set upload path', path)
    uploader.setParams({ dir: path })
}