// JavaScript Document

function getQueryVariable(variable)
{
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i=0;i<vars.length;i++) {
        var pair = vars[i].split("=");
        if(pair[0] === variable){return pair[1];}
    }
    return(false);
}

function iniParam() {
    var laypage = layui.laypage;

    //页面效果
    $("#keyWord").focus(function () {
        $(this).parent().addClass("search-border");
    }).blur(function () {
        $(this).parent().removeClass("search-border");
    }).keydown(function (e) {
        if (e.which == 13) { //监听回车事件
            search($(this).val());
            return false;
        }
    });

    //搜索
    $('#search').click(function () {
        var value = $("#keyWord").val();
        search(value);
    });

    function search(value) {
        if (!value) {
            layer.tips('关键字都没输入,你想搜啥...', '#keyWord', { tips: [1, '#659FFD'] });
            $("#keyWord")[0].focus(); //使文本框获得焦点
            return;
        } else if (value.length > 20) {
            layer.tips('关键字最长不能超过20...', '#keyWord', { tips: [1, '#659FFD'] });
            $("#keyWord")[0].focus(); //使文本框获得焦点
            return;
        }

        window.location.href = "/blog/list?item=" + value
    }

    // 获取url参数
    var offset = getQueryVariable("offset")
    var limit = getQueryVariable("limit")
    var type = getQueryVariable("type")
    var item = getQueryVariable("item")

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
		// theme: "page",
		jump: function (obj, first) {
			if (!first) { //首次不执行
			    var urlValue = "/blog/list?offset=" + obj.curr + "&limit=" + obj.limit
                if (type !== false) {
                    urlValue = urlValue + "&type=" + type
                }
                if (item !== false) {
                    urlValue = urlValue + "&item=" + item
                }
			    window.location.href = urlValue
			}
		}
	});
}
