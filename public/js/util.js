/**
 * 工具类，需要JQuery支持！！
 * by Liu Han
 * 2020-8-16
 */

function map_to_obj(map) {
    let obj = Object.create(null);
    for (let [k, v] of map) {
        obj[k] = v;
    }
    return obj
}

function http_get(url, success_callback, error_callback) {
    $.ajax({
        url: url,
        type: 'get',
        success: function(data) {
            success_callback(data)
        },
        error: function(e) {
            error_callback(e)
        }
    });
}

function http_request(method, url, json_data, success_callback = null, error_callback = null) {
    $.ajax({
        url: url,
        type: method,
        dataType: 'json',
        contentType: "application/json", //json 格式作为提交
        data: JSON.stringify(json_data),
        success: function(data) {
            if (success_callback != null) {
                success_callback(data)
            }
        },
        error: function(e) {
            if (error_callback != null) {
                error_callback(e.statusText)
            }
        }
    });
}

function http_post(url, json_data, success_callback, error_callback) {
    http_request('post', url, json_data, success_callback, error_callback)
}

function http_delete(url, json_data, success_callback, error_callback) {
    http_request('delete', url, json_data, success_callback, error_callback)
}

function http_put(url, json_data, success_callback, error_callback) {
    http_request('put', url, json_data, success_callback, error_callback)
}

function log(txt) {
    console.log(txt)
}

function get_suffix(filename) {
    var index = filename.lastIndexOf(".");
    if (index != -1)
        return filename.substring(index + 1, filename.length).toLowerCase();
    else return "";
}