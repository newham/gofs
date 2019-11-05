var isAll = false;
/*user display:none*/
// $('#btn-delete').hide();

function selectAll() {
    if (isAll) {
        $('input:checkbox').each(function () {
            this.checked = false;
            // $(this).hide();
        });
        setBtn(false);
    } else {
        $('input:checkbox').each(function () {
            this.checked = true;
            // $(this).show();
        });
        setBtn(true);
    }
    isAll = !isAll;
}

function deleteAll() {
    var files = "";
    $('input:checkbox').each(function () {
        if (this.checked) {
            files = files + "|" + $(this).val();
        }

    });
    if (files == "") {
        alert("please select one");
        return;
    }
    if (confirm("Ensure to delete[" + files + "]?")) {
        $.post(
            "/del?type=array",
            {array: files},
            function (data) {
                $('body').html(data);
            }
        )
    }

}

$('input:checkbox').click(function () {
    var isChecked = false;
    $('input:checkbox').each(function () {
        if (this.checked) {
            isChecked = true;
            isAll = true;
        }
    });
    setBtn(isChecked);
})

function setBtn(isChecked) {
    if (isChecked) {
        $('#btn-all').html('<span class="glyphicon glyphicon-remove" aria-hidden="true"></span>');
        // $('#btn-all').removeClass("btn-success");
        // $('#btn-all').addClass("btn-warning");
        $('#btn-delete').show();
    } else {
        $('#btn-all').html('<span class="glyphicon glyphicon-ok" aria-hidden="true"></span>');
        // $('#btn-all').addClass("btn-success");
        // $('#btn-all').removeClass("btn-warning");
        $('#btn-delete').hide();
    }

}

var isGird = true;

function share(filename) {
    $.get(
        "/share?name="+filename,
        function (data) {
            // alert();
            url ="http://"+window.location.host +"/download?shareKey="+ data.shareKey;

            // $('#shareFileUrl').attr("href",url); 
            $('#shareFileName').html(data.file);
            $('#shareFileUrl').val(url);
            $('#shareModal').modal("show");
            
        }
    )
}