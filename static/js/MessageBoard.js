// JavaScript Document

function iniParam() {
    var form = layui.form,laypage = layui.laypage,layedit = layui.layedit;
	
    //评论和留言的编辑器
	for(var i=1;i<9;i++){
		layedit.build('demo-'+i.toString(), {
			height: 150,
			tool: ['face', '|', 'link'],
		});
	}
	

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

	// 获取url参数
	var offset = getQueryVariable("offset")
	var limit = getQueryVariable("limit")

	if (offset === false) {
		offset = 1
	}

	if (limit === false) {
		limit = 5
	}

	var total_count = $("#page").attr("data-count")

	laypage.render({
		elem: 'page',
		count: total_count, //数据总数通过服务端得到
		limit: limit, //每页显示的条数。laypage将会借助 count 和 limit 计算出分页数。
		curr: offset, //起始页。一般用于刷新类型的跳页以及HASH跳页。
		first: '首页',
		last: '尾页',
		layout: ['prev', 'page', 'next', 'skip'],
		//theme: "page",
		jump: function (obj, first) {
			if (!first) { //首次不执行
				window.location.href = "/blog/leave?offset=" + obj.curr + "&limit=" + obj.limit
			}
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

function aricleComment(aricleId) {
	var commentValue = $(".blog-leave-comment").val()
	if (commentValue !== "") {
		var jsonData = {
			article_id: aricleId,
			type: 3,  // 留言
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

function addReply(replyId, commentId, aimsUserId, idKey) {
	var replyValue = $("#" + idKey + "-" + replyId).val()
	if (replyValue !== "") {
		var jsonData = {
			comment_id: commentId,
			type: 2,
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


