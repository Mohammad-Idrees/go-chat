package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"project/config"
	db "project/db/sqlc"
	"project/logger"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	MEMBERSHIP_CHANNEL = "membership"
)

type Hub struct {
	redisClient          *redis.Client
	ChannelSubscriptions map[int64]int           // map[channelId]no of client subscribers
	ChannelPubSub        map[int64]*redis.PubSub // map[channel id]pub sub object
	Membership           map[int64]*db.Membership
	Clients              map[int64]*Client
	AddClient            chan *Client
	RemoveClient         chan *Client
	ReadBroadcast        chan string
	WriteBroadcast       chan *Message
	wg                   *sync.WaitGroup
	MembershipUpdates    chan *db.Membership
	serverName           string
}

func newHub(wg *sync.WaitGroup, cfg *config.StartupConfig, redisClient *redis.Client) *Hub {
	return &Hub{
		redisClient:          redisClient,
		ChannelSubscriptions: make(map[int64]int),
		ChannelPubSub:        make(map[int64]*redis.PubSub),
		Membership:           make(map[int64]*db.Membership),
		Clients:              map[int64]*Client{},
		AddClient:            make(chan *Client, 10),
		RemoveClient:         make(chan *Client, 10),
		ReadBroadcast:        make(chan string, 10),
		WriteBroadcast:       make(chan *Message, 10),
		wg:                   wg,
		MembershipUpdates:    make(chan *db.Membership, 10),
		serverName:           cfg.Server.Name,
	}
}

func InitHub(wg *sync.WaitGroup, cfg *config.StartupConfig, redisClient *redis.Client) *Hub {
	hub := newHub(wg, cfg, redisClient)
	pubsub := hub.redisClient.Subscribe(context.Background(), MEMBERSHIP_CHANNEL)
	go hub.run()
	go hub.membershipUpdatesReader(pubsub)
	return hub
}

func (hub *Hub) membershipUpdatesReader(pubsub *redis.PubSub) {
	hub.wg.Add(1)
	defer hub.wg.Done()
	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			logger.Error(context.Background(), "membershipUpdatesReader", logger.Field("redis receive message error", err.Error()))
			continue
		}

		membership := &db.Membership{}
		err = json.Unmarshal([]byte(msg.Payload), membership)
		if err != nil {
			logger.Error(context.Background(), "ReadPump", logger.Field("unmarshal error", err.Error()))
			continue
		}

		// membership client not connected to this hub, continue
		if _, ok := hub.Clients[membership.UserID]; !ok {
			continue
		}

		// new membership identified for the client connected to the hub.
		hub.Clients[membership.UserID].Memberships = append(hub.Clients[membership.UserID].Memberships, membership)
		hub.addSubscription(membership)
		hub.WriteBroadcast <- &Message{
			Content:   "user joined the channel",
			ChannelId: membership.ChannelID,
			Username:  hub.Clients[membership.UserID].Username,
		}
	}
}

func (h *Hub) run() {
	h.wg.Add(1)
	defer h.wg.Done()

	for {
		select {
		case client := <-h.AddClient:
			h.addClient(client)

		case client := <-h.RemoveClient:
			h.removeClient(client)

		case message := <-h.ReadBroadcast:
			h.readBroadcast(message)

		case message := <-h.WriteBroadcast:
			h.writeBroadcast(message)

		case membership := <-h.MembershipUpdates:
			h.membershipUpdates(membership)
		}
	}
}

func (hub *Hub) membershipUpdates(membership *db.Membership) {

	membershipBytes, errr := json.Marshal(membership)
	if errr != nil {
		logger.Error(context.Background(), "membershipUpdates", logger.Field("marshal error", errr.Error()))
		return
	}

	membershipStr := string(membershipBytes)
	err := hub.redisClient.Publish(context.Background(), MEMBERSHIP_CHANNEL, membershipStr)
	if err.Err() != nil {
		logger.Error(context.Background(), "membershipUpdates", logger.Field("redis publish error", err.Err().Error()))
		return
	}
}

func (hub *Hub) addClient(client *Client) {
	// client exists
	if _, ok := hub.Clients[client.Id]; ok {
		return
	}

	// add client in memory
	hub.Clients[client.Id] = client
	for _, membership := range client.Memberships {
		hub.WriteBroadcast <- &Message{
			Content:   fmt.Sprintf("user connected to server %s", hub.serverName),
			ChannelId: membership.ChannelID,
			Username:  hub.Clients[membership.UserID].Username,
		}
		hub.addSubscription(membership)
	}
}

func (hub *Hub) removeClient(client *Client) {
	// client does not exist
	if _, ok := hub.Clients[client.Id]; !ok {
		return
	}

	// remove client from memory
	for _, membership := range client.Memberships {
		hub.WriteBroadcast <- &Message{
			Content:   fmt.Sprintf("user disconnected from server %s", hub.serverName),
			ChannelId: membership.ChannelID,
			Username:  hub.Clients[membership.UserID].Username,
		}
		hub.removeSubscription(membership)
	}
	delete(hub.Clients, client.Id)
}

func (hub *Hub) addSubscription(membership *db.Membership) {
	// add subscription to inmemory
	if _, ok := hub.Membership[membership.ID]; !ok {
		hub.Membership[membership.ID] = membership
	}

	// redis channel subscription exists
	if _, ok := hub.ChannelSubscriptions[membership.ChannelID]; ok {
		subscribers := hub.ChannelSubscriptions[membership.ChannelID]
		hub.ChannelSubscriptions[membership.ChannelID] = subscribers + 1
		return
	}

	channelId := fmt.Sprint(membership.ChannelID)
	pubsub := hub.redisClient.Subscribe(context.Background(), channelId)
	hub.ChannelPubSub[membership.ChannelID] = pubsub
	go hub.startSubscription(pubsub)
	hub.ChannelSubscriptions[membership.ChannelID] = 1
}

func (hub *Hub) startSubscription(pubsub *redis.PubSub) {
	hub.wg.Add(1)
	defer hub.wg.Done()

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			logger.Error(context.Background(), "startSubscription", logger.Field("redis receive message error", err.Error()))
			continue
		}
		hub.ReadBroadcast <- msg.Payload
	}

}

func (hub *Hub) removeSubscription(membership *db.Membership) {
	// remove subscription to inmemory
	delete(hub.Membership, membership.ID)

	// redis channel subscription does not exists
	if _, ok := hub.ChannelSubscriptions[membership.ChannelID]; !ok {
		return
	}

	subscribers := hub.ChannelSubscriptions[membership.ChannelID]
	if subscribers == 1 {
		delete(hub.ChannelPubSub, membership.ChannelID)
		delete(hub.ChannelSubscriptions, membership.ChannelID)
	}
}

func (hub *Hub) readBroadcast(msgStr string) {
	msg := &Message{}
	err := json.Unmarshal([]byte(msgStr), msg)
	if err != nil {
		logger.Error(context.Background(), "readBroadcast", logger.Field("unmarshal error", err.Error()))
		return
	}

	channelId := msg.ChannelId
	for _, membership := range hub.Membership {
		if membership.ChannelID == channelId {
			hub.Clients[membership.UserID].MessageChan <- msg
		}
	}
}

func (hub *Hub) writeBroadcast(msg *Message) {

	msgBytes, errr := json.Marshal(msg)
	if errr != nil {
		logger.Error(context.Background(), "writeBroadcast", logger.Field("marshal error", errr.Error()))
		return
	}

	msgStr := string(msgBytes)
	channelId := fmt.Sprint(msg.ChannelId)
	err := hub.redisClient.Publish(context.Background(), channelId, msgStr)
	if err.Err() != nil {
		logger.Error(context.Background(), "writeBroadcast", logger.Field("redis publish error", err.Err().Error()))
		return
	}

	// hub.ReadBroadcast <- msgStr
}
