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

func GetDocumentTree(array []blog.Tree, buf *bytes.Buffer, args ...interface{}) {
	detailsId := args[0].(int)
	detailsMap := args[1].(map[string]blog.Tree)

	buf.WriteString("<ul>")

	for index, item := range array {
		if detailsId != 0 && detailsId == item.Id {
			// 获取上一篇文章ID
			if index != 0 {
				detailsMap["onContent"] = array[index-1]
			}
			if index+1 < len(array) {
				detailsMap["underContent"] = array[index+1]
			}
		}

		buf.WriteString("<li>")
		if detailsId == item.Id {
			buf.WriteString(fmt.Sprintf("<a href='/blog/series_details?series_id=%v&details_id=%v' style='color: #1E9FFF' title='%s'>%s</a>", args[2], item.Id, item.Text, item.Text))
		} else {
			buf.WriteString(fmt.Sprintf("<a href='/blog/series_details?series_id=%v&details_id=%v' title='%s'>%s</a>", args[2], item.Id, item.Text, item.Text))
		}

		if len(item.Nodes) > 0 {
			GetDocumentTree(item.Nodes, buf, args...)
		}
		buf.WriteString("</li>")
	}
	buf.WriteString("</ul>")
}
