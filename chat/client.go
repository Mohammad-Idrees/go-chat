// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"context"
	"encoding/json"
	db "project/db/sqlc"
	"project/logger"

	"github.com/gorilla/websocket"
)

type Message struct {
	Content   string `json:"content"`
	ChannelId int64  `json:"channelId"`
	Username  string `json:"username"`
}

type Client struct {
	Id          int64
	Username    string
	Conn        *websocket.Conn
	MessageChan chan *Message
	Memberships []*db.Membership
}

func (c *Client) WritePump(hub *Hub) {
	hub.wg.Add(1)
	defer hub.wg.Done()

	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.MessageChan
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadPump(hub *Hub) {
	hub.wg.Add(1)
	defer hub.wg.Done()

	defer func() {
		hub.RemoveClient <- c
		c.Conn.Close()
	}()

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error(context.Background(), "ReadPump", logger.Field("error", err.Error()))
			}
			break
		}

		msg := &Message{}
		err = json.Unmarshal(messageBytes, msg)
		if err != nil {
			logger.Error(context.Background(), "ReadPump", logger.Field("unmarshal error", err.Error()))
			continue
		}
		if msg.Username != c.Username {
			logger.Error(context.Background(), "ReadPump", logger.Field("error", "unauthorized user"))
			continue
		}

		hub.WriteBroadcast <- msg
	}
}
