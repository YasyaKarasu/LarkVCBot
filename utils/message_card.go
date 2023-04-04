package utils

import "github.com/YasyaKarasu/feishuapi"

func DefaultMarkdownMessageCard(template feishuapi.MessageCardTitleTemplate, title string, content string) string {
	card, _ := feishuapi.NewMessageCard().
		WithConfig(
			feishuapi.NewMessageCardConfig().
				WithEnableForward(true).
				WithUpdateMulti(true).
				Build(),
		).
		WithHeader(
			feishuapi.NewMessageCardHeader().
				WithTemplate(template).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent(title).
						Build(),
				).
				Build(),
		).
		WithElements([]feishuapi.MessageCardElement{
			feishuapi.NewMessageCardMarkdown().
				WithContent(content).
				Build(),
		}).Build().String()
	return card
}

func DefaultMarkdownMessageCardInfo(title string, content string) string {
	return DefaultMarkdownMessageCard(feishuapi.TemplateBlue, title, content)
}

func DefaultMarkdownMessageCardWarn(title string, content string) string {
	return DefaultMarkdownMessageCard(feishuapi.TemplateOrange, title, content)
}

func DefaultMarkdownMessageCardError(title string, content string) string {
	return DefaultMarkdownMessageCard(feishuapi.TemplateRed, title, content)
}
