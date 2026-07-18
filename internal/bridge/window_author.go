package bridge

import "github.com/pkg/browser"

const footerAuthorHomeURL = "https://github.com/Liu-fu-gui/cursor-api-bridge"

var footerAuthorInfo = FooterAuthorInfo{
	ButtonText:        "GitHub",
	DialogTitle:       "Cursor API Bridge",
	DialogContent:     "本项目以 MIT 许可证开源。源代码、版本更新与问题反馈请访问 GitHub 项目页。",
	DialogConfirmText: "打开 GitHub",
	DialogCancelText:  "关闭",
}

// FooterAuthorInfo 定义首页底部作者入口的展示信息。
type FooterAuthorInfo struct {
	ButtonText        string `json:"buttonText"`
	DialogTitle       string `json:"dialogTitle"`
	DialogContent     string `json:"dialogContent"`
	DialogConfirmText string `json:"dialogConfirmText"`
	DialogCancelText  string `json:"dialogCancelText"`
}

// GetFooterAuthorInfo 返回首页底部作者入口的展示信息。
func (s *WindowService) GetFooterAuthorInfo() FooterAuthorInfo {
	return footerAuthorInfo
}

// OpenFooterAuthorHome 打开作者主页。
func (s *WindowService) OpenFooterAuthorHome() error {
	return browser.OpenURL(footerAuthorHomeURL)
}
