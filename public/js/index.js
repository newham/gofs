function set_data(data) {
    set_file_list(data.Files)
    set_path(data.PathArray)
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
    // console.log(btn_group_show)
    if (btn_group_show == 1) {
        $('#btn-group-1').hide()
        $('#btn-group-2').show()
        btn_group_show = 2
    } else {
        console.log('checked_map')
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

$('#btn-upload').click((e) => {
    // console.log('btn-upload')
    if (is_drag_show) {
        $("#upload-box").hide()
        $('#btn-upload').removeClass('bg-red')
    } else {
        $("#upload-box").show()
        $('#btn-upload').addClass('bg-red')
    }
    is_drag_show = !is_drag_show
})

function switch_class(id, c) {
    var item = $(`#${id}`)
    if (item.hasClass(c)) {
        item.removeClass(c)
    } else {
        item.addClass(c)
    }
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

$("#input-new-folder").blur(function() {
    $("#new-folder").hide()
});

function new_folder() {
    $("#new-folder").hide()
    console.log('new folder')
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

$('#input-new-folder').keydown(function(e) {
    if (e.keyCode == 13) { //inter
        new_folder()
        return false
    } else if (e.keyCode == 27) { //esc
        $("#new-folder").hide()
    }
});

function getFileName(file) {
    var pos = file.lastIndexOf('/')
    return file.substring(pos + 1)
}

$(document).ready(() => {
    get_data('').then((data) => {
        set_data(data)
        set_upload(data.Path, () => {
            get_data_http(data.Path).then((data2) => {
                set_data(data2)
            })
        })
    })
})