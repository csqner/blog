/*
  @Author : lanyulei
*/

package errno

var (
	// success
	OK = &Errno{Code: 10000, Message: "OK"}

	// common errors
	InternalServerError = &Errno{Code: 10001, Message: "内部服务器错误"}
	ErrBind             = &Errno{Code: 10002, Message: "将请求Body绑定到Struct发送错误"}
	ErrTokenInvalid     = &Errno{Code: 10004, Message: "Token无效"}

	// blog
	ErrSelectList    = &Errno{Code: 20001, Message: "获取文章列表失败"}
	ErrCreateBlog    = &Errno{Code: 20002, Message: "创建文章失败"}
	ErrSelectDetails = &Errno{Code: 20003, Message: "获取文章详情失败"}
	ErrUpdateView    = &Errno{Code: 20004, Message: "更新文件访问量失败"}
	ErrTotalCount    = &Errno{Code: 20005, Message: "获取文章总数失败"}
	ErrTotalView     = &Errno{Code: 20006, Message: "获取文章总阅读量失败"}
	ErrLastTime      = &Errno{Code: 20007, Message: "获取文章最近更新时间失败"}
	ErrAwesomeCount  = &Errno{Code: 20008, Message: "获取总赞数失败"}
	ErrUpdateAwesome = &Errno{Code: 20008, Message: "更新赞数失败"}

	// blog tag
	ErrCreateTag     = &Errno{Code: 30001, Message: "创建标签失败"}
	ErrSelectTag     = &Errno{Code: 30002, Message: "查询标签失败"}
	ErrCreateBlogTag = &Errno{Code: 30003, Message: "关联文章标签失败"}

	// blog type
	ErrSelectTypeList = &Errno{Code: 40001, Message: "获取分类列表失败"}
	ErrTypeCount      = &Errno{Code: 40004, Message: "获取分类总数失败"}

	// blog links
	ErrSelectLinks = &Errno{Code: 50001, Message: "获取友情链接失败"}

	// blog user
	ErrSelectUser = &Errno{Code: 60001, Message: "获取用户信息失败"}
	ErrNotLogin   = &Errno{Code: 60002, Message: "未登陆"}

	// blog comment
	ErrCreateComment = &Errno{Code: 70001, Message: "发表评论失败"}
	ErrSelectComment = &Errno{Code: 70002, Message: "查看评论信息失败"}
)
