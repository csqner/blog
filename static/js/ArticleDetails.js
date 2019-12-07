// JavaScript Document

function iniParam() {
    var form = layui.form,laypage = layui.laypage,layedit = layui.layedit;

 	layer.photos({
		photos: '#details-content',
		anim: 5 //0-6的选择，指定弹出图片动画类型，默认随机（请注意，3.0之前的版本用shift参数）
	});

	//评论和留言的编辑器
	for(var i=0;i<9;i++){
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
	
	laypage.render({
		elem: 'page',
		count: 10, //数据总数通过服务端得到
		limit: 5, //每页显示的条数。laypage将会借助 count 和 limit 计算出分页数。
		curr: 1,
		first: '首页',
		last: '尾页',
		layout: ['prev', 'page', 'next', 'skip'],
		//theme: "page",
		jump: function (obj, first) {
			if (!first) { //首次不执行
				layer.msg("第"+obj.curr+"页");

			}
		}
	});
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
	var commentValue = $(".blog-aricle-comment").val()
	if (commentValue !== "") {
		var jsonData = {
			article_id: aricleId,
			type: 1,
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

function addAwesome(content_id) {
	$.get("/blog/add-awesome", {
		content_id: content_id
	}, function (response) {
		if (response.code === 60002) {
			layer.msg("评论需登陆，请点击右下角哆啦A梦登陆")
		} else if (response.code === 10000) {
			window.location.reload()
		} else {
			layer.msg(response.message)
		}
	})
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
