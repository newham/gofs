function set_user() {
    var user_name = get_user_name()
    $('#user-name').html(user_name)
    $('#btn-user').html(user_name)
}

function get_user_name() {
    return $.cookie('user')
}

function set_path_cookie(path) {
    $.cookie('path', path)
}

function get_path_cookie() {
    if (!$.cookie('path')) {
        return "/"
    }
    return $.cookie('path')
}