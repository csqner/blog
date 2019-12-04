package tree

import (
	"blog/models/blog"
	"bytes"
	"fmt"
)

/*
  @Author : lanyulei
*/

//获取所有权限
//递归实现(返回树状结果得数据)
func SuperSeriesTree(allSeries []blog.SeriesDetailsStruct, pid int) []blog.Tree {
	var treeArray []blog.Tree
	for _, v := range allSeries {
		if pid == v.Parent {
			treeValue := blog.Tree{}
			treeValue.Id = v.Id
			treeValue.Pid = v.Parent
			treeValue.Text = v.Title
			seriesChildren := SuperSeriesTree(allSeries, v.Id)
			treeValue.Nodes = seriesChildren
			treeArray = append(treeArray, treeValue)
		}
	}
	return treeArray
}

func GetDocumentTree(array []blog.Tree, buf *bytes.Buffer) {
	buf.WriteString("<ul>")

	for _, item := range array {
		buf.WriteString("<li>")
		buf.WriteString(fmt.Sprintf("<a href='http://baidu.com' title='%s'>%s</a>", item.Text, item.Text))
		if len(item.Nodes) > 0 {
			GetDocumentTree(item.Nodes, buf)
		}
		buf.WriteString("</li>")
	}
	buf.WriteString("</ul>")
}
