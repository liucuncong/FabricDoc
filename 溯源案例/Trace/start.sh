#!/bin/bash
echo "启动容器: orderer..."
docker-compose -f docker-orderer.yaml up -d
echo "启动容器: dairy.peer0 and cli ..."
docker-compose -f docker-peer0-dairy.yaml up -d
echo "启动容器: dairy.peer1 ..."
docker-compose -f docker-peer1-dairy.yaml up -d
echo "启动容器: process.peer0 ..."
docker-compose -f docker-peer0-process.yaml up -d
echo "启动容器: process.peer1 ..."
docker-compose -f docker-peer1-process.yaml up -d
echo "启动容器: sell.peer0 ..."
docker-compose -f docker-peer0-sell.yaml up -d
echo "启动容器: sell.peer1 ..."
docker-compose -f docker-peer1-sell.yaml up -d
