package delivery

import (
	"net/http"
	"project/chat"
	db "project/db/sqlc"
	"project/models/request"
	"project/service"
	"project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WSHandler struct {
	hub     *chat.Hub
	userSvc service.UserService
	repo    db.Repository
}

func NewWSHandler(hub *chat.Hub, userSvc service.UserService, repo db.Repository) *WSHandler {
	return &WSHandler{
		hub:     hub,
		userSvc: userSvc,
		repo:    repo,
	}
}

func ConfigureWSHandler(router *gin.RouterGroup, hub *chat.Hub, userSvc service.UserService, repo db.Repository) {
	wsHandler := NewWSHandler(hub, userSvc, repo)
	addWSHandlerRoutes(router, wsHandler)
}

func addWSHandlerRoutes(router *gin.RouterGroup, wsHandler *WSHandler) {
	router.GET("/channels", wsHandler.GetChannels)
	router.GET("/memberships", wsHandler.GetMemberships)
	router.POST("/channels", wsHandler.CreateChannel)
	router.GET("/ws/join", wsHandler.JoinChat)
	router.GET("/channels/join/:channelId", wsHandler.JoinChannel)
}

func (h *WSHandler) GetChannels(c *gin.Context) {
	ctx := c.Request.Context()
	channels, err := h.repo.GetChannels(ctx)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channels)
}

func (h *WSHandler) GetMemberships(c *gin.Context) {
	ctx := c.Request.Context()
	memberships, err := h.repo.GetMemberships(ctx)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, memberships)
}

func (h *WSHandler) CreateChannel(c *gin.Context) {
	ctx := c.Request.Context()
	var createChannelRequest request.CreateChannelRequest
	if err := c.ShouldBindJSON(&createChannelRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.repo.CreateChannel(ctx, createChannelRequest.Name)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channel)
}

func (h *WSHandler) JoinChat(c *gin.Context) {
	ctx := c.Request.Context()
	// var joinChatRequest request.JoinChatRequest
	// if err := c.ShouldBindUri(&joinChatRequest); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	userId, err := strconv.ParseInt(c.Query("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId query param"})
		return
	}

	getUserByIdRequest := request.GetUserByIdRequest{
		Id: userId,
	}
	user, err := h.userSvc.GetUserById(ctx, &getUserByIdRequest)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	memberships, err := h.repo.GetMembershipsByUserId(ctx, user.Id)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// upgrade to websocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	username := user.Username
	client := &chat.Client{
		Id:          user.Id,
		Username:    username,
		Conn:        conn,
		MessageChan: make(chan *chat.Message, 10),
		Memberships: memberships,
	}

	// add client to server
	h.hub.AddClient <- client

	// write & read loops
	go client.WritePump(h.hub)
	go client.ReadPump(h.hub)
}

func (h *WSHandler) JoinChannel(c *gin.Context) {
	ctx := c.Request.Context()
	var joinChannelRequest request.JoinChannelRequest
	if err := c.ShouldBindUri(&joinChannelRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindQuery(&joinChannelRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getUserByIdRequest := request.GetUserByIdRequest{
		Id: joinChannelRequest.UserId,
	}
	user, err := h.userSvc.GetUserById(ctx, &getUserByIdRequest)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.repo.GetChannelById(ctx, joinChannelRequest.ChannelId)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	createMembershipParams := &db.CreateMembershipParams{
		UserID:    user.Id,
		ChannelID: channel.ID,
	}

	membership, err := h.repo.CreateMembership(ctx, createMembershipParams)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	h.hub.MembershipUpdates <- membership
	c.JSON(http.StatusOK, membership)
}
