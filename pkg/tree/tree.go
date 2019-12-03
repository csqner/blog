package tree

import "blog/models/blog"

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
