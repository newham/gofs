function set_data(data) {
    btn_group_show = 1
    checked_map = new Map()
    checked_count = 0
    is_drag_show = false
    set_file_list(data.Files)
    set_path(data.PathArray)
    show_btn_group(1)
    show_select_info()
}

function set_path(pathArray) {
    // 清空原来的
    $("#path-array").html('')

    // console.log(pathArray.length, pathArray)
    if (pathArray.length == 1) {
        var item = pathArray[0]
        $("#path-array").append(`<li class="item all cur"><a href="/folder/${item[1]}" title="${item[0]}">${item[0]}</a></li>`)
        return
    }
    for (i = 0; i < pathArray.length; i++) {
        var item = pathArray[i]
        if (i == 0) {
            $("#path-array").append(`<li class="item all"><a href="/folder/${item[1]}"><i class="icon"></i>${item[0]}</a></li>`)
        } else if (i == pathArray.length - 1) {
            $("#path-array").append(`<li class="item cur"><i class="icon icon-bread-next"></i> <a href="/folder/${item[1]}" title="${item[0]}">${item[0]}</a></li>`)
        } else {
            $("#path-array").append(`<li class="item"><i class="icon icon-bread-next"></i> <a href="/folder/${item[1]}" title="${item[0]}">${item[0]}</a></li>`)
        }
    }
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
                                    <div class="info"><a href="javascript:void(0)" title="新建文件夹" class="tit" style="display: none;">新建文件夹</a>
                                        <span class="fileedit"><input type="text" class="ui-input" id="input-new-folder" onkeydown="new_folder_keydown()" onblur="hide_new_folder()"></span>
                                    </div>
                                </div>
                                <div class="item-info"><span class="item-info-list"><span class="txt txt-time"></span></span> <span class="item-info-list"><span class="txt txt-size">-</span></span>
                                    <!---->
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
        $("#file-list").append(getFileLi(-1, 'file', `href="/folder/${data.PathArray[index][1]}"`, "../上级目录", "", ""))
    }
    // var action = `href="/folder/${data.PathArray[index][1]}"`
    // var f = {Name:data.PathArray[index][0],ModTime:"",Size:""}

    // $("#file-list").append(file_li)
    for (i = 0; i < file_list.length; i++) {
        var f = data.Files[i]
        var icon = f.Type
        var action = `href="/file/${f.Path}" target="blank"`
        if (f.Type == 'folder') {
            icon = "file"
            action = `href="/folder/${f.Path}"`
        } else if (f.Type == 'pic') {
            action = `href="javascript:void(0)" onclick="preview(${i})"`
        } else if (f.Type == 'flv' || f.Type == 'video') {
            action = `href="/video/${f.Path}" target="blank"`
        } else if (f.Type == 'audio') {
            action = `href="javascript:void(0)" onclick="play_audio('${f.Path}')"`
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
                                        <div class="info">
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

var btn_group_show = 1

var checked_map = new Map()

var checked_count = 0

function check_file(i) {
    console.log('check:', i)
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

$('#btn-user').click((e) => {
    switch_class('menu-user', 'act')
})

var is_drag_show = false

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

function map_to_json(map) {
    let obj = Object.create(null);
    for (let [k, v] of map) {
        obj[k] = v;
    }
    return JSON.stringify(obj)
}

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

    //2.开始创建
    $.ajax({
        url: '/folder/',
        type: 'put',
        dataType: 'json',
        contentType: "application/json", //form 格式作为提交
        data: JSON.stringify({
            "dir": data.Path + new_folder_name //千万别掉了data.Path
        }),
        success: function(d) {
            console.log('new folder success')
            reload_page()
        },
        error: function(e) {
            console.log(e)
        }
    });
}

$("#btn-close-audio").click((e) => {
    $("#audio").attr('src', '')
    $("#player-audio").hide()
})

function play_audio(path) {
    console.log(encodeURI(`/file/${path}`))
    $("#audio").attr('src', encodeURI(`/file/${path}`))
    $("#audio-title").html(getFileName(path))
    $("#player-audio").show()
        // $("#audio").play()
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

function delete_file() {
    $.ajax({
        url: '/file/',
        type: 'delete',
        dataType: 'json',
        contentType: "application/json", //json 格式作为提交
        data: map_to_json(checked_map),
        success: function(d) {
            console.log('delete success')
            reload_page()
        },
        error: function(e) {
            console.log(e)
        }
    });
}


$("#btn-delete").click((e) => {
    console.log('delete')
    delete_file()
})

// $("#app").click((e) => {
//     if(event.currentTarget.id in ['','','']){
//         return
//     }
//     // hide_menu()
// })

$("#_layout_main").click((e) => {
    console.log('hide')
    hide_menu()
})

$(".layout-aside").click((e) => {
    console.log('hide')
    hide_menu()
})


$("#btn-select-all").click((e) => {
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


// $("#toolbar").click((e) => {
//     console.log('hide')
//     hide_menu()
// })

function hide_menu() {
    show_upload_box(false)
    remove_class('menu-new', 'act')
    remove_class('menu-user', 'act')
}

function reload_page() {
    get_data_http(data.Path).then((data_new) => {
        // 设置新数据
        set_data(data_new)
    })
}

$(document).keyup(function(e) {
    var key = e.which || e.keyCode;
    if (key == 27) { //esc
        show_upload_box(false)
        hide_menu()
    }
});

$(document).ready(() => {
    get_data('').then((data) => {
        // 设置数据
        set_data(data)

        //设置下载页面
        set_upload(data.Path, () => {
            //这是回调函数
            get_data_http(data.Path).then((data_new) => {
                set_data(data_new)
            })
        })
    })
})