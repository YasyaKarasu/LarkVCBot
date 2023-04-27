package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/global"

	"github.com/YasyaKarasu/feishuapi"
)

func groupHelp(messageevent *chat.MessageEvent, args ...any) {
	card, _ := feishuapi.NewMessageCard().
		WithConfig(
			feishuapi.NewMessageCardConfig().
				WithEnableForward(true).
				WithUpdateMulti(true).
				Build(),
		).
		WithHeader(
			feishuapi.NewMessageCardHeader().
				WithTemplate(feishuapi.TemplateBlue).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent("💡 使用说明").
						Build(),
				).
				Build(),
		).
		WithElements([]feishuapi.MessageCardElement{
			feishuapi.NewMessageCardMarkdown().
				WithContent(
					"1. 第一次使用时，向机器人发送 *初始化* 指令，机器人会首先向指令发送人**私聊**发送鉴权请求，" +
						"鉴权完成后，机器人将在**群聊**内提示输入知识库ID（获取方式见下图）。" +
						"完成上述操作后，机器人将自动进行初始化，并在完成后提示。" +
						"需要注意的是，由于上述操作需要涉及将机器人添加为知识库管理员，因此操作者需要是**知识库管理员**。",
				).
				Build(),
			feishuapi.NewMessageCardImage().
				WithAlt(
					feishuapi.NewMessageCardPlainText().
						WithContent("打开知识空间设置").
						Build(),
				).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent("打开知识空间设置").
						Build(),
				).
				WithMode(feishuapi.ModeFitHorizontal).
				WithPreview(true).
				WithImageKey("img_v2_a29dd47f-68c7-43dd-a39e-d5916ac41b9g").
				Build(),
			feishuapi.NewMessageCardImage().
				WithAlt(
					feishuapi.NewMessageCardPlainText().
						WithContent("获取知识库ID").
						Build(),
				).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent("获取知识库ID").
						Build(),
				).
				WithMode(feishuapi.ModeFitHorizontal).
				WithPreview(true).
				WithImageKey("img_v2_70e3eb9c-9a2b-4f20-a397-4fd44f5b50dg").
				Build(),
			feishuapi.NewMessageCardMarkdown().
				WithContent(
					"2. 机器人完成初始化后将为群聊创建一张日历，当机器人检测到**此日历下**有新日程，" +
						"将自动将日程添加到知识库日程排期中，并提示前往排期页面设置会议主持人，以接受会前统计信息。" +
						"同时，在日程前一天，机器人将在知识库创建会议记录文档。记得从日程发起视频会议（点日程开始前飞书弹出的发起视频会议）哦~",
				),
			feishuapi.NewMessageCardImage().
				WithAlt(
					feishuapi.NewMessageCardPlainText().
						WithContent("PC端选择日历").
						Build(),
				).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent("PC端选择日历").
						Build(),
				).
				WithMode(feishuapi.ModeFitHorizontal).
				WithPreview(true).
				WithImageKey("img_v2_c14a5068-3b73-4aa8-8105-1956ba870f8g").
				Build(),
			feishuapi.NewMessageCardImage().
				WithAlt(
					feishuapi.NewMessageCardPlainText().
						WithContent("移动端选择日历").
						Build(),
				).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent("移动端选择日历").
						Build(),
				).
				WithMode(feishuapi.ModeFitHorizontal).
				WithPreview(true).
				WithImageKey("img_v2_07d0f924-5d41-41f6-891b-78d977a1dfcg").
				Build(),
		}).Build().String()
	global.FeishuClient.MessageSend(
		feishuapi.GroupChatId,
		messageevent.Message.Chat_id,
		feishuapi.Interactive,
		card,
	)
}

func p2pHelp(messageevent *chat.MessageEvent, args ...any) {
	card, _ := feishuapi.NewMessageCard().
		WithConfig(
			feishuapi.NewMessageCardConfig().
				WithEnableForward(true).
				WithUpdateMulti(true).
				Build(),
		).
		WithHeader(
			feishuapi.NewMessageCardHeader().
				WithTemplate(feishuapi.TemplateBlue).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent("💡 使用说明").
						Build(),
				).
				Build(),
		).
		WithElements([]feishuapi.MessageCardElement{
			feishuapi.NewMessageCardMarkdown().
				WithContent("会议机器人对私聊暂时没有主动功能哦～").
				Build(),
		}).Build().String()
	global.FeishuClient.MessageSend(
		feishuapi.UserOpenId,
		messageevent.Sender.Sender_id.Open_id,
		feishuapi.Interactive,
		card,
	)
}
