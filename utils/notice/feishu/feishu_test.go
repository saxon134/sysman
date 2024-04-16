package feishu

import (
	"testing"
)

func TestFeishu(t *testing.T) {
	var client = New("https://open.cn/open-apis/bot/v2/hook/xxx", "")
	var eles = make([]*CardItem, 0, 2)
	eles = append(eles, &CardItem{
		Tag:  "div",
		Text: CardText{Tag: "lark_md", Content: "**标重点：**好好学习，天天向上"},
	}, &CardItem{
		Tag: "hr",
	}, &CardItem{
		Tag:     "markdown",
		Content: "- [ ] 支持以 PDF 格式导出文稿\n- [ ] 改进 Cmd 渲染算法\n- [x] 新增 Todo 列表功能\n- [x] 修复 LaTex 公式渲染问题",
	}, &CardItem{
		Tag: "hr",
	}, &CardItem{
		Tag:  "div",
		Text: CardText{Tag: "lark_md", Content: "**还是要好好学习，天天向上**"},
	}, &CardItem{
		Tag:     "markdown",
		Content: "```gantt\n    title 项目开发流程\n    section 项目确定\n        需求分析       :a1, 2016-06-22, 3d\n        可行性报告     :after a1, 5d\n        概念验证       : 5d\n    section 项目实施\n        概要设计      :2016-07-05  , 5d\n        详细设计      :2016-07-08, 10d\n        编码          :2016-07-15, 10d\n        测试          :2016-07-22, 5d\n    section 发布验收\n        发布: 2d\n        验收: 3d\n```",
	})
	client.SendCard("", eles)
}
