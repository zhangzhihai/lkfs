package lkfs

//搜索模块
//每存入一片文字内容，提交id到搜索模块，先进行分词，生成的关键词进行索引
//索引设计为，关键词对应的最新的文章索引ID，最新的ID对应上一篇文章id
//关键词库为map[string]int32(文件位置) 索引为 int32,int32 文章id,上一篇文章id
//关键词库的初始化与设置,生成一个大的分词结构体
/*
 1. 读取词库的初始化生成map to byte
 2. 关键词库文件与关键词搜索文件分开来存放
*/
