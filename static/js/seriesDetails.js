function iniParam() {
    $(".btn-reply").click(function(){
        if ($(this).text() === '回复') {
            $(this).parent().next().removeClass("layui-hide");
            $(this).text('回复');
            $(this).text('收起');
        }
        else {
            $(this).parent().next().addClass("layui-hide");
            $(this).text('回复');
        }
    });
}

function getQueryVariable(variable) {
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i=0;i<vars.length;i++) {
        var pair = vars[i].split("=");
        if(pair[0] === variable){return pair[1];}
    }
    return(false);
}

function searchTitleValue() {
    var series_id = getQueryVariable("series_id")
    if (series_id !== false) {
        window.location.href = "/blog/series_details?series_id= "+ series_id +" &searchValue=" + $("#searchInputValue").val();
    }
}

function getBrowser() {
    var explorer =navigator.userAgent ;
    var browser = "";
    if (explorer.indexOf("MSIE") >= 0) {
        browser = "IE"
    } else if (explorer.indexOf("Firefox") >= 0) {
        browser = "Firefox"
    } else if(explorer.indexOf("Chrome") >= 0){
        browser = "Chrome"
    } else if(explorer.indexOf("Opera") >= 0){
        browser = "Opera"
    } else if(explorer.indexOf("Safari") >= 0){
        browser = "Safari"
    } else if(explorer.indexOf("Netscape")>= 0) {
        browser = "Netscape"
    } else {
        browser = "其他浏览器"
    }
    return browser
}

function addReply(replyId, commentId, aimsUserId, idKey) {
    var replyValue = $("#" + idKey + "-" + replyId).val()
    if (replyValue !== "") {
        var jsonData = {
            comment_id: commentId,
            type: 1,
            aims_user_id: aimsUserId,
            content: replyValue,
            browser: getBrowser()
        }
        $.post("/blog/save-reply", jsonData, function (response) {
            if (response.code === 60002) {
                layer.msg("评论需登陆，请点击右下角哆啦A梦登陆")
            } else if (response.code === 10000) {
                window.location.reload()
            } else {
                layer.msg(response.message)
            }
        })
    } else {
        layer.msg("回复内容不能为空！");
    }
}

function aricleComment(aricleId) {
    var commentValue = $(".blog-aricle-comment").val()
    if (commentValue !== "") {
        var jsonData = {
            article_id: aricleId,
            type: 2,
            content: commentValue,
            browser: getBrowser()
        }

        $.post("/blog/save-comment", jsonData, function (response) {
            console.log(response)
            if (response.code === 60002) {
                layer.msg("评论需登陆，请点击右下角哆啦A梦登陆")
            } else if (response.code === 10000) {
                window.location.reload()
            } else {
                layer.msg(response.message)
            }
        })
    } else {
        layer.msg("评论内容不能为空！");
    }
}