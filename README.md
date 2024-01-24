# Real-Time Chat Application with Redis Pub/Sub and PostgreSQL

## Overview

This Golang-based chat application provides real-time communication between users through WebSocket technology. The application leverages Redis Pub/Sub for seamless communication across multiple WebSocket servers, ensuring scalability and responsiveness. PostgreSQL is used for persistent storage of chat messages.

## Features

- Real-time chat between multiple users.
- Multiple webservers using docker-compose
- NGINX as a load balancer with least_conn as balancing algorithm.
- Scalable architecture using Redis Pub/Sub for inter-server communication.
- WebSocket for bidirectional communication between clients and servers.
- PostgreSQL for persistent storage of user, different channel information.


https://github.com/Mohammad-Idrees/go-chat/assets/64984896/ad351498-c248-4d11-9bb4-543059403275



