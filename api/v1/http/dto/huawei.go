package dto

// ClickAction 点击行为
type ClickAction struct {
	Url    string
	Action Action // 点击行为类型
}

// Action 枚举表示点击行为类型
type Action int

const (
	// ActionTypeOpenCustomPage 表示打开应用自定义页面
	ActionTypeOpenCustomPage Action = iota + 1

	// ActionTypeOpenURL 表示点击后打开特定URL
	ActionTypeOpenURL

	// ActionTypeOpenApp 表示点击后打开应用
	ActionTypeOpenApp
)

type HuaweiPushRequestData struct {
	DryRun      bool
	Foreground  bool
	TTL         string
	Type        string
	Title       string
	Body        string
	Category    string
	Icon        string
	Img         string
	Sound       string
	ClickAction ClickAction // 点击行为
}

type ApnsPushRequestData struct {
}
