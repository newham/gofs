var data
var btn_group_show = 1
var checked_map = new Map()
var checked_count = 0
var is_drag_show = false

// 以下数据不随页面更新而更新
var audio_map = new Map()
var audio_map_len = 0
var audio_index = 0

const default_rank_by_type = ['folder', 'docx', 'doc', 'xlsx', 'xls', 'pptx', 'ppt', 'pdf', 'txt', 'md', 'jpg', 'jpeg', 'png', 'bmp', 'gif', 'mp3', 'mp4', 'flv', 'mkv', 'avi', 'zip', '7z', 'tar', 'rar', 'gz', 'exe']

function get_data_http(path = '', isBase64 = false) {
    return new Promise(function(resolve, reject) {
        http_post('/folder/', {
            "dir": path,
            "base64": isBase64,
        }, (d) => {
            data = d
            resolve(d);
        }, (e) => {
            reject(e);
        })
    })
}

function set_data(d) {
    //绑定数据
    data = d

    //排序,默认是根据类型分组，文件夹优先
    data.Files = sort_by_type(data.Files)

    //初始化参数
    btn_group_show = 1
    checked_map = new Map()
    checked_count = 0
    is_drag_show = false

    //初始化界面
    set_file_list(data.Files)
    set_path(data.PathArray)
    if (isCut) {
        show_btn_group(2)
    } else {
        show_btn_group(1)
    }
    show_select_info()
}

function set_path(pathArray) {
    // 清空原来的
    $("#path-array").html('')

    // console.log(pathArray.length, pathArray)
    if (pathArray.length == 1) {
        var item = pathArray[0]
        $("#path-array").append(`<li class="item all cur"><a href="javascript:void(0)" onclick="load_page('${item[1]}',true)" title="${item[0]}">${item[0]}</a></li>`)
        return
    }
    for (i = 0; i < pathArray.length; i++) {
        var item = pathArray[i]
        if (i == 0) {
            $("#path-array").append(`<li class="item all"><a href="javascript:void(0)" onclick="load_page('${item[1]}',true)"><i class="icon"></i>${item[0]}</a></li>`)
        } else if (i == pathArray.length - 1) {
            $("#path-array").append(`<li class="item cur"><i class="icon icon-bread-next"></i><a href="javascript:void(0)" onclick="load_page('${item[1]}',true)" title="${item[0]}">${item[0]}</a></li>`)
        } else {
            $("#path-array").append(`<li class="item"><i class="icon icon-bread-next"></i> <a href="javascript:void(0)" onclick="load_page('${item[1]}',true)" title="${item[0]}">${item[0]}</a></li>`)
        }
    }
}

//根据文件后缀分类，时间复杂度2n，空间复杂度n
function sort_by_type(file_list, rank = default_rank_by_type) {
    //文件名后缀
    var typeMap = new Map()
    console.log('sort_by_type')
    for (i in file_list) {
        var file = file_list[i]
        if (!typeMap.get(file.Suffix)) {
            typeMap.set(file.Suffix, [file])
        } else {
            var array = typeMap.get(file.Suffix)
            array.push(file)
            typeMap.set(file.Suffix, array)
        }

    }
    file_list = [] //清空未排序的列表
    for (var j in rank) { //1.先把rank里的写入
        var k1 = rank[j]
        if (typeMap.has(k1)) {
            file_list.push.apply(file_list, typeMap.get(k1)) //数组拼接
            typeMap.delete(k1) //2.从map中删除
        }
    }
    //2.把剩下的随机写入
    for (var k2 of typeMap.keys()) {
        file_list.push.apply(file_list, typeMap.get(k2))
    }
    return file_list
}

function set_file_list(file_list) {
    // 先清空原来的
    $("#file-list").html('')

    var new_folder_li = `<li class="list-group-item checked" id="new-folder" style="display: none;">
                            <div class="item-inner">
                                <div class="item-tit">
                                    <div class="label"><i class="icon icon-check-s icon-checkbox"></i></div>
                                    <div class="thumb">
                                        <i class="icon icon-m icon-file-m"></i>
                                    </div>
                                    <div class="info">
                                        <a href="javascript:void(0)" title="新建文件夹" class="tit" style="display: none;">新建文件夹</a>
                                        <span class="fileedit">
                                            <input type="text" class="ui-input" id="input-new-folder" onkeydown="new_folder_keydown()" onblur="hide_new_folder()">
                                        </span>
                                    </div>
                                </div>
                                <div class="item-info">
                                    <span class="item-info-list"><span class="txt txt-time"></span></span> <span class="item-info-list"><span class="txt txt-size">-</span></span>
                                </div>
                            </div>
                        </li>`

    // 新建文件夹li
    $("#file-list").append(new_folder_li)

    // var icon ="file"
    // var index = 0
    // 当存在上级目录时，显示上级目录
    if (data.PathArray.length > 1) {
        index = data.PathArray.length - 2
        $("#file-list").append(getFileLi(-1, 'file', `href="javascript:void(0)" onclick="load_page('${data.PathArray[index][1]}',true)"`, "../上级目录", "", ""))
    }
    // var action = `href="/folder/${data.PathArray[index][1]}"`
    // var f = {Name:data.PathArray[index][0],ModTime:"",Size:""}

    // $("#file-list").append(file_li)
    for (i = 0; i < file_list.length; i++) {
        var f = file_list[i] //这里不是 data.Files!!!请注意，否则排序无效
        var icon = f.Type
        var action = `href="/file/${f.Path}"`
        if (f.Type == 'folder') {
            icon = "file"
                // console.log("path", `${f.Path}`)
            action = `href="javascript:void(0)" onclick="load_page('${f.Path}',true)"` // action = `href="/folder/${f.Path}"`
        } else if (f.Type == 'pic') {
            action = `href="javascript:void(0)" onclick="preview(${i})"`
        } else if (f.Type == 'flv' || f.Type == 'video') {
            action = `href="/video/${f.Path}" target="blank"`
        } else if (f.Type == 'audio') {
            action = `href="javascript:void(0)" onclick="play_audio('${f.Path}','${f.Name}')"`
        } else if (f.Editable) {
            action = `href="/edit/${f.Path}" target="blank"`
        }
        $("#file-list").append(getFileLi(i, icon, action, f.Name, f.ModTime, f.Size))
    }
    $("#count").html(file_list.length)
}

function getFileLi(i, icon, action, name, mod_time, size) {
    var check = `<i class="icon icon-check-s icon-checkbox" onclick="check_file(${i})"></i>`
    if (i == -1) {
        check = '<i class="icon"></i>'
    }
    var file_li = `<li class="list-group-item checked" id="check-${i}">
                                <div class="item-inner">
                                    <div class="item-tit">
                                        <div class="label">
                                            ${check}
                                        </div>
                                        <div class="thumb">
                                            <i class="icon icon-m icon-${icon}-m"></i></div>
                                        <div class="info" id="file-info-${i}">
                                            <a ${action} title="${name}" class="tit">${name}</a><br>
                                        </div>
                                    </div>
                                    <div class="item-info"><span class="item-info-list">
                                        <span class="txt txt-time">${mod_time}</span></span> 
                                        <span class="item-info-list"><span class="txt txt-size">${size}</span></span>
                                    </div>
                                </div>
                            </li>`
    return file_li
}

$("#btn-rename").click((e) => {
    if (checked_count == 0) {
        return
    }
    show_rename()
})

var temp_old_file_html = ''
var temp_old_file_index = 0

function show_rename() {
    // log('show_rename')
    for (var i of checked_map.keys()) {
        temp_old_file_index = i
        var btn = $(`#file-info-${i}`)
        temp_old_file_html = btn.html()
        temp_old_file_name = data.Files[i].Name
        btn.html(`<span class="fileedit"><input type="text" value="${temp_old_file_name}" class="ui-input" id="input-rename" onkeydown="rename_keydown()" onblur="hide_rename()"></span>`)
        $("#input-rename").focus()
    }
}

function hide_rename() {
    // log('hide_rename', temp_old_file_index)
    $(`#file-info-${temp_old_file_index}`).html(temp_old_file_html)
}

function rename_keydown() {
    if (event.keyCode == 13) { //inter
        rename()
        return false
    } else if (event.keyCode == 27) { //esc
        hide_rename()
    }
}

function rename() {
    log('do rename')
    http_post('/file/rename', { old: data.Files[temp_old_file_index].Name, new: $('#input-rename').val(), path: data.Files[temp_old_file_index].Path }, reload_page, (e) => {
        alert("重命名失败：", e)
    })
}

function check_file(i) {
    //如果有还未粘贴的文件，不允许点击选择
    if (isCut) {
        alert('请先粘贴，再选择！')
        log('check pass')
        return
    }

    // console.log('check:', i)
    switch_class(`check-${i}`, 'act')
    if ($(`#check-${i}`).hasClass('act')) {
        // console.log(data.Files[i].Path)
        checked_map.set(i, data.Files[i].Path)
        checked_count++
    } else {
        // console.log('delete',checked_map.get(i))
        checked_map.delete(i)
        checked_count--
    }
    show_select_info()
        // console.log(btn_group_show)
    show_btn_group()

    //如果checked_count > 1则隐藏重命名按钮
    if (checked_count > 1) {
        $("#btn-rename").hide();
    } else {
        $("#btn-rename").show();
    }

}

function show_btn_group(id) {
    if (id != null) {
        if (id == 1) {
            btn_group_show = 2
            checked_count = 0
        } else {
            btn_group_show = 1
        }
    }
    if (btn_group_show == 1) {
        $('#btn-group-1').hide()
        $('#btn-group-2').show()
        btn_group_show = 2

        //hide drag
        show_upload_box(false)
    } else {
        // console.log('checked_map')
        if (checked_count == 0) {
            $('#btn-group-2').hide()
            $('#btn-group-1').show()
            btn_group_show = 1
        }
    }
}

$('#btn-new').click((e) => {
    switch_class('menu-new', 'act')
})

$('#btn-download-list').click((e) => {
    if (downloaded_count > 0) {
        switch_class('menu-download', 'act')
    }
})

$('#btn-user').click((e) => {
    switch_class('menu-user', 'act')
})

$('#formFileInputCt').click((e) => {
    // console.log('btn-upload')
    show_upload_box()
})

$('#task-box').click((e) => {
    show_upload_box()
})

function show_upload_box(is_show) {
    if (is_show != null) {
        is_drag_show = !is_show
    }
    if (is_drag_show) {
        $("#upload-box").hide()
        $('#btn-upload').removeClass('bg-red')
        $('#formFileInputCt').removeClass('bg-red')

        if (fileCount > 0) {
            $('#task-box').show()
        }
    } else {
        $("#upload-box").show()
        $('#btn-upload').addClass('bg-red')
        $('#formFileInputCt').addClass('bg-red')

        $('#task-box').hide()
    }
    is_drag_show = !is_drag_show
}

function switch_class(id, c) {
    var item = $(`#${id}`)
    if (item.hasClass(c)) {
        item.removeClass(c)
    } else {
        item.addClass(c)
    }
}

function remove_class(id, c) {
    var item = $(`#${id}`)
    item.removeClass(c)
}

var current = 0

function preview(id) {
    current = id
    set_img(id)
    $("#img-preview").show()
}

function set_img(id) {
    var img = data.Files[id]
    $("#img-title").html(img.Name)
    $("#img").attr('src', encodeURI('/file/' + img.Path))
    console.log(img.Path)
}

$('#btn-close-preview').click((e) => {
    $("#img-preview").hide()
})

function img_to(step) {
    while (current < data.Files.length) {
        if (step > 0 && current == data.Files.length - 1) {
            // 从头开始
            current = -1
        }
        if (step < 0 && current == 0) {
            current = data.Files.length
        }
        current += step
        if (data.Files[current].Type == 'pic') {
            set_img(current)
            break
        }

    }
}

$('#btn-img-next').click((e) => {
    img_to(+1)
})

$('#btn-img-pre').click((e) => {
    img_to(-1)
})

$("#btn-new-folder").click((e) => {
    $("#new-folder").show()
    $("#input-new-folder").focus()
})

// $("#input-new-folder").blur(function() {
//     $("#new-folder").hide()
// });

function hide_new_folder() {
    $("#new-folder").hide()
}

function new_folder() {
    var new_folder_name = $("#input-new-folder").val()

    //1.检查是否重名
    for (i = 0; i < data.Files.length; i++) {
        var file = data.Files[i]
        if (file.Name == new_folder_name) {
            alert("文件夹已经存在！")
            return
        }
    }
    console.log("new folder:", data.Path + new_folder_name)

    //2.开始创建
    http_put('/folder/', {
        "dir": data.Path + new_folder_name,
    }, reload_page, console.log)
}

$("#btn-close-audio").click((e) => {
    $("#audio").attr('src', '')
    $("#player-audio").hide()

    //clean data
    audio_map = new Map()
    audio_map_len = 0
    audio_index = 0
})

function play_audio(path, name) {
    console.log('play', name)
    $("#audio").attr('src', encodeURI(`/file/${path}`))
    $("#audio-title").html(name)
    $("#player-audio").show()

    // $("#audio").play()

    set_audio_list(name)
}

function set_audio_list(name) {
    for (i = 0; i < data.Files.length; i++) {
        var f = data.Files[i]
        if (f.Type == 'audio') {
            audio_map.set(audio_map_len, f)
            if (name == f.Name) {
                audio_index = audio_map_len
            }
            audio_map_len++
        }
    }
}

$("#audio").on('ended', (e) => {
    console.log('audio ended')
    audio_next()
})

function audio_back() {
    if (--audio_index == -1) {
        audio_index = audio_map_len - 1
    }
    var f = audio_map.get(audio_index)
    play_audio(f.Path, f.Name)
}

function audio_next() {
    if (++audio_index == audio_map_len) {
        audio_index = 0
    }
    var f = audio_map.get(audio_index)
    play_audio(f.Path, f.Name)
}

// $('#input-new-folder').keydown(function(e) {

// });

function new_folder_keydown() {
    if (event.keyCode == 13) { //inter
        new_folder()
        return false
    } else if (event.keyCode == 27) { //esc
        $("#new-folder").hide()
    }
}

function getFileName(file) {
    var pos = file.lastIndexOf('/')
    return file.substring(pos + 1)
}

function download_files() {
    log("download files")
    http_post('/download', map_to_obj(checked_map), (data) => {
        log(data.tmp)
        for (var i of checked_map.keys()) {
            check_file(i)
        }
        //添加到完成列表
        add_downloaded_list(data.tmp)
    }, console.log)
}

var downloaded_count = 0

function add_downloaded_list(tmp) {
    var name = "..." + tmp.substring(tmp.length - 28, tmp.length)
    var li = `<li class="menu-item" id="btn-new-folder">
                <span class="txt">
                    <a href="/tmp/${tmp}" target="blank">${name}</a>
                </span>
            </li>`

    $("#downloaded-count").html(++downloaded_count)
    $("#downloaded-list").append(li)

    //显示或隐藏任务按钮
    if (downloaded_count > 0) {
        $("#menu-download").show()
    } else {
        $("#menu-download").hide()
    }
}

function delete_file() {
    http_delete('/file/', map_to_obj(checked_map), reload_page, console.log)
}

$("#btn-delete").click((e) => {
    if (checked_count == 0) {
        return
    }
    console.log('delete')
    delete_file()
})

$("#btn-download").click((e) => {
    if (checked_count == 0) {
        return
    }
    console.log('download')
    download_files()
})

$("#_layout_main").click((e) => {
    console.log('hide')
    hide_menu()
})

$(".layout-aside").click((e) => {
    console.log('hide')
    hide_menu()
})


$("#btn-select-all").click((e) => {
    if (isCut) {
        alert('请先粘贴，再选择！')
        log('check pass')
    }
    console.log('select all')
    for (i = 0; i < data.Files.length; i++) {
        check_file(i)
    }
})

function show_select_info() {
    $("#selected-count").html(checked_count)
    if (checked_count > 0) {
        $("#select-all-box").removeClass('up')
        $("#select-all-box").addClass('checked')
        $("#select-all-box").addClass('cur')
    } else {
        $("#select-all-box").removeClass('checked')
        $("#select-all-box").removeClass('cur')
        $("#select-all-box").addClass('up')
    }

    if (checked_count == data.Files.length) {
        $("#select-all-item").addClass('act')
    } else {
        $("#select-all-item").removeClass('act')
    }

}

function hide_menu() {
    show_upload_box(false)
    remove_class('menu-new', 'act')
    remove_class('menu-user', 'act')
    remove_class('menu-download', 'act')
    remove_class('menu-more', 'act')
}

var isCut = false
var tmp_checked_map

function cut(ifCut) {
    log('cut', ifCut)
    isCut = ifCut
    if (ifCut) {
        tmp_checked_map = checked_map
    } else {
        tmp_checked_map = null
    }
    show_paste()
    reload_page()
}

function show_paste() {
    if (isCut) {
        $("#btn-paste").show()
        $("#btn-cut").hide()
        $("#btn-cut-cancel").show()
    } else {
        $("#btn-paste").hide()
        $("#btn-cut").show()
        $("#btn-cut-cancel").hide()
    }
}

function paste() {
    http_post("/file/move", { checkedMap: map_to_obj(tmp_checked_map), dir: data.Path }, (d) => {
        isCut = false
        show_paste()
        reload_page()
    }, (e) => {
        alert("粘贴失败！文件已经存在/路径无效")
    })
}

function reload_page() {
    load_page(data.Path)
}

function get_current() {
    return data.Path
}

function load_page(path, isBase64 = false) {
    get_data_http(path, isBase64).then((data) => {
        // 设置新数据
        set_data(data)

        //设置uploader 的 path
        set_upload_path(data.Path)
    })
}

$(document).keyup(function(e) {
    var key = e.which || e.keyCode;
    if (key == 27) { //esc
        show_upload_box(false)
        hide_menu()
    }
});