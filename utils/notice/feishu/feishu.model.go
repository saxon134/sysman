package feishu

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// PostItem 富文本内容项，参考：https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
type PostItem struct {
	Tag    string `json:"tag"` //text,a,at
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

// CardItem 卡片内容项，参考：https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-structure/card-content
type CardItem struct {
	Tag     string   `json:"tag"` //div,hr,button,markdown,note,hr
	Content string   `json:"content"`
	Text    CardText `json:"text,omitempty"`
}

type CardText struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}
