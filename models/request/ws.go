package request

type CreateChannelRequest struct {
	Name string `json:"name" binding:"required"`
}

type JoinChatRequest struct {
	UserId int64 `form:"userId"`
}

type JoinChannelRequest struct {
	UserId    int64 `form:"userId"`
	ChannelId int64 `uri:"channelId"`
}
