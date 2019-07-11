> 我们可以自己组建一个Fabric网路, 网络结构如下: 
>
> - 排序节点 1 个
> - 组织个数 2 个, 分别为go和cpp, 每个组织分别有两个peer节点, 用户个数为3

| 机构名称 | 组织标识符 |  组织ID   |
| :------: | :--------: | :-------: |
|  Go学科  |   org_go   | OrgGoMSP  |
|   CPP    |  org_cpp   | OrgCppMSP |

**一些理论基础:**

- 域名
  - baidu.com
  - jd.com
  - taobao.com
- msp（可以理解为一个账号）
  - Membership service provider (MSP)是一个提供虚拟成员操作的管理框架的组件。
  - 账号
    - 都谁有msp
      - 每个节点都有一个msp账号
      - 每个用户都有msp账号
- 锚节点
  - 代表所属组织和其他组织进行通信的节点

- Peer节点的分类

  ![1542205733340](D:/%E9%BB%91%E9%A9%AC%E5%8C%BA%E5%9D%97%E9%93%BE/%E8%B5%84%E6%96%99/%E7%AC%AC%E4%B8%83%E9%98%B6%E6%AE%B5/hyperledger/day03-%E9%80%9A%E9%81%93%E6%93%8D%E4%BD%9C%E5%92%8C%E6%99%BA%E8%83%BD%E5%90%88%E7%BA%A6/01-%E6%95%99%E5%AD%A6%E8%B5%84%E6%96%99/assets/1542205733340.png)

## 1. 生成fabric证书（cryptogen）

### 1.1 命令介绍

```shell
$cryptogen --help
```

### 1.2 证书的文件的生成 - yaml

- **配置文件的模板**

  ```yaml
  # ---------------------------------------------------------------------------
  # "OrdererOrgs" - Definition of organizations managing orderer nodes
  # ---------------------------------------------------------------------------
  OrdererOrgs:	# 排序节点组织信息，这个名字不能改（一般不会只有一个，需要对order节点集群，协同工作，进行负载均衡）
    # ---------------------------------------------------------------------------
    # Orderer
    # ---------------------------------------------------------------------------
    - Name: Orderer	# 排序节点组织的名字（可以改）
      Domain: example.com	# 根域名, 排序节点组织的根域名（ip的别名，真是的生产环境，这个域名必须要去注册）
      Specs:
        - Hostname: orderer # 访问这台orderer节点对应的域名为: orderer.example.com
        - Hostname: order2 # 访问这台orderer节点对应的域名为: order2.example.com
  # ---------------------------------------------------------------------------
  # "PeerOrgs" - Definition of organizations managing peer nodes
  # ---------------------------------------------------------------------------
  PeerOrgs:
    # ---------------------------------------------------------------------------
    # Org1
    # ---------------------------------------------------------------------------
    - Name: Org1	# 第一个组织的名字, 自己指定
      Domain: org1.example.com	# 访问第一个组织用到的根域名
      EnableNodeOUs: true			# 链码编写时是否支持node.js
      Template:					# 模板, 根据默认的规则生成2个peer存储数据的节点
        Count: 2 # 1. peer0.org1.example.com 2. peer1.org1.example.com
      Users:	   # 创建的普通用户的个数（操作节点的用户，根据项目进行评估得出），还会默认生成一个管理员用户
        Count: 3
        
    # ---------------------------------------------------------------------------
    # Org2: See "Org1" for full specification
    # ---------------------------------------------------------------------------
    - Name: Org2
      Domain: org2.example.com
      EnableNodeOUs: true
      Template:
        Count: 2
      Specs:
        - Hostname: hello
      Users:
        Count: 1
  ```

  > **上边使用的域名, 在真实的生成环境中需要注册备案（买一个域名，自己的服务器有一个ip地址，让别人给你绑定，自己绑定不了）, 测试环境, 域名自己随便指定就可以**

- 根据要求编写好的配置文件, 配置文件名: crypto-config.yaml（起这个名字一看就知道给cryptogen命令用的）

  ```yaml
  # crypto-config.yaml
  # ---------------------------------------------------------------------------
  # "OrdererOrgs" - Definition of organizations managing orderer nodes
  # ---------------------------------------------------------------------------
  OrdererOrgs:
    # ---------------------------------------------------------------------------
    # Orderer
    # ---------------------------------------------------------------------------
    - Name: Orderer
      Domain: itcast.com
      Specs:
        - Hostname: orderer
  
  # ---------------------------------------------------------------------------
  # "PeerOrgs" - Definition of organizations managing peer nodes
  # ---------------------------------------------------------------------------
  PeerOrgs:
    # ---------------------------------------------------------------------------
    # Org1
    # ---------------------------------------------------------------------------
    - Name: OrgGo
      Domain: orggo.itcast.com
      EnableNodeOUs: true
      Template:
        Count: 2
      Users:
        Count: 3
  
    # ---------------------------------------------------------------------------
    # Org2: See "Org1" for full specification
    # ---------------------------------------------------------------------------
    - Name: OrgCpp
      Domain: orgcpp.itcast.com
      EnableNodeOUs: true
      Template:
        Count: 2
      Specs:
        - Hostname: hello
      Users:
        Count: 3
  
  ```

- 通过命令生成证书文件

  ```shell
  $ cryptogen generate --config=crypto-config.yaml
  ```

Specs:

- Hostname: orderer

 Template:
      Count: 2

用这两个都行，区别是Specs可以指定域名。也可以两个混用，如下（3个节点）：

- Name: OrgCpp
    Domain: orgcpp.itcast.com
    EnableNodeOUs: true
    Template:
      Count: 2
    Specs:
    
      - Hostname: hello
    
      Users:
        Count: 3

## 2. 创始块文件和通道文件的生成（configtxgen）

### 2.1 命令介绍

```shell
$ configtxgen --help 
  # 输出创始块区块文件的路径和名字
  `-outputBlock string`
  # 指定创建的channel的名字, 如果没指定系统会提供一个默认的名字.
  `-channelID string`
  # 表示输通道文件路径和名字
  `-outputCreateChannelTx string`
  # 指定配置文件中的节点
  `-profile string`
  # 更新channel的配置信息
  `-outputAnchorPeersUpdate string`
  # 指定所属的组织名称
  `-asOrg string`
  # 要想执行这个命令, 需要一个配置文件 configtx.yaml
```

**要想执行这个命令, 需要一个配置文件 configtx.yaml**（必须叫这个名字）

### 2.2 创始块/通道文件的生成

- **配置文件的编写** - <font color="red">参考模板</font>

  ```yaml
  
  ---
  ################################################################################
  #
  #   Section: Organizations
  #
  #   - This section defines the different organizational identities which will
  #   be referenced later in the configuration.
  #
  ################################################################################
  Organizations:			# 固定的不能改
      - &OrdererOrg		# 排序节点组织, 自己起个名字
          Name: OrdererOrg	# 排序节点的组织名
          ID: OrdererMSP		# 排序节点组织的ID
          MSPDir: crypto-config/ordererOrganizations/example.com/msp # 组织的msp账号信息（在前面创建的crypto-config.yaml文件里）
  
      - &Org1			# 第一个组织, 名字自己起
          Name: Org1MSP # 第一个组织的名字
          ID: Org1MSP		# 第一个组织的ID
          MSPDir: crypto-config/peerOrganizations/org1.example.com/msp
          AnchorPeers: # 锚节点
              - Host: peer0.org1.example.com  # 指定一个peer节点的域名
                Port: 7051					# 端口不要改（节点启动后是一个容器，7051是容器的这个进程监听的端口，发数据要发到容器的这个端口）
  
      - &Org2
          Name: Org2MSP
          ID: Org2MSP
          MSPDir: crypto-config/peerOrganizations/org2.example.com/msp
          AnchorPeers:
              - Host: peer0.org2.example.com
                Port: 7051
  
  ################################################################################
  #
  #   SECTION: Capabilities, 在fabric1.1之前没有, 设置的时候全部设置为true
  #   
  ################################################################################
  Capabilities:
      Global: &ChannelCapabilities
          V1_1: true
      Orderer: &OrdererCapabilities
          V1_1: true
      Application: &ApplicationCapabilities
          V1_2: true
  
  ################################################################################
  #
  #   SECTION: Application
  #
  ################################################################################
  Application: &ApplicationDefaults
      Organizations:
  
  ################################################################################
  #
  #   SECTION: Orderer
  #
  ################################################################################
  Orderer: &OrdererDefaults
      # Available types are "solo" and "kafka"
      # 共识机制 == 排序算法
      OrdererType: solo	# 排序方式（solo就一个节点），kafka需要集群
      Addresses:			# orderer节点的地址
          - orderer.example.com:7050	# 端口不要改（接收数据的端口）
  
  	# BatchTimeout,MaxMessageCount,AbsoluteMaxBytes只要一个满足, 区块就会产生
      BatchTimeout: 2s	# 多长时间产生一个区块
      BatchSize:
          MaxMessageCount: 10		# 交易的最大数据量, 数量达到之后会产生区块, 建议100左右
          AbsoluteMaxBytes: 99 MB # 数据量达到这个值, 会产生一个区块, 32M/64M
          PreferredMaxBytes: 512 KB  #不需要改
      Kafka: #上面是solo，这里就没有什么意义了
          Brokers: #指Kafka的服务器
              - 127.0.0.1:9092
      Organizations:
  
  ################################################################################
  #
  #   Profile
  #
  ################################################################################
  Profiles:	# 不能改
      TwoOrgsOrdererGenesis:	# 区块名字, 随便改
          Capabilities:
              <<: *ChannelCapabilities
          Orderer:
              <<: *OrdererDefaults
              Organizations:
                  - *OrdererOrg
              Capabilities:
                  <<: *OrdererCapabilities
          Consortiums:
              SampleConsortium:	# 这个名字可以改
                  Organizations:
                      - *Org1
                      - *Org2
      TwoOrgsChannel:	# 通道名字, 可以改
          Consortium: SampleConsortium	# 这个名字对应93行
          Application:
              <<: *ApplicationDefaults
              Organizations:
                  - *Org1
                  - *Org2
              Capabilities:
                  <<: *ApplicationCapabilities
  
  ```

  - &Org1			# 第一个组织, 名字自己起
          Name: Org1MSP # 第一个组织的名字
          ID: Org1MSP		# 第一个组织的ID
          MSPDir: crypto-config/peerOrganizations/org1.example.com/msp
          AnchorPeers: # 锚节点
              - Host: peer0.org1.example.com  # 指定一个peer节点的域名
                Port: 7051					# 端口不要改

 Organizations:

- *Org1

*Org1代表的就是&Org1这个整体

锚节点负责peer节点组织之间通信的，leader节点是peer节点组织和order组织节点进行通信的

- 按照要求编写的配置文件

  ```yaml
  # configtx.yaml
  ---
  ################################################################################
  #
  #   Section: Organizations
  #
  ################################################################################
  Organizations:
      - &OrdererOrg
          Name: OrdererOrg
          ID: OrdererMSP
          MSPDir: crypto-config/ordererOrganizations/itcast.com/msp
  
      - &org_go
          Name: OrgGoMSP
          ID: OrgGoMSP
          MSPDir: crypto-config/peerOrganizations/orggo.itcast.com/msp
          AnchorPeers:
              - Host: peer0.orggo.itcast.com
                Port: 7051
  
      - &org_cpp
          Name: OrgCppMSP
          ID: OrgCppMSP
          MSPDir: crypto-config/peerOrganizations/orgcpp.itcast.com/msp
          AnchorPeers:
              - Host: peer0.orgcpp.itcast.com
                Port: 7051
  
  ################################################################################
  #
  #   SECTION: Capabilities
  #
  ################################################################################
  Capabilities:
      Global: &ChannelCapabilities
          V1_1: true
      Orderer: &OrdererCapabilities
          V1_1: true
      Application: &ApplicationCapabilities
          V1_2: true
  
  ################################################################################
  #
  #   SECTION: Application
  #
  ################################################################################
  Application: &ApplicationDefaults
      Organizations:
  
  ################################################################################
  #
  #   SECTION: Orderer
  #
  ################################################################################
  Orderer: &OrdererDefaults
      # Available types are "solo" and "kafka"
      OrdererType: solo
      Addresses:
          - orderer.itcast.com:7050
      BatchTimeout: 2s
      BatchSize:
          MaxMessageCount: 100
          AbsoluteMaxBytes: 32 MB
          PreferredMaxBytes: 512 KB
      Kafka:
          Brokers:
              - 127.0.0.1:9092
      Organizations:
  
  ################################################################################
  #
  #   Profile
  #
  ################################################################################
  Profiles:
      ItcastOrgsOrdererGenesis:
          Capabilities:
              <<: *ChannelCapabilities
          Orderer:
              <<: *OrdererDefaults
              Organizations:
                  - *OrdererOrg
              Capabilities:
                  <<: *OrdererCapabilities
          Consortiums:
              SampleConsortium:
                  Organizations:
                      - *org_go
                      - *org_cpp
      ItcastOrgsChannel:
          Consortium: SampleConsortium
          Application:
              <<: *ApplicationDefaults
              Organizations:
                  - *org_go
                  - *org_cpp
              Capabilities:
                  <<: *ApplicationCapabilities
  
  ```

- **执行命令生成文件**

  > <font color="red">-profile 后边的参数从configtx.yaml中的Profiles 里边的配置项</font>

  - 生成创始块文件

    ```shell
    $ configtxgen -profile ItcastOrgsOrdererGenesis -outputBlock ./genesis.block
    - 在当前目录下得到一个文件: genesis.block
    ```

  - 生成通道文件

    ```shell
    $ configtxgen -profile ItcastOrgsChannel -outputCreateChannelTx channel.tx -channelID itcastchannel
    #后面的节点加入通道都是要通过通道ID来完成的
  ```
  
- 生成锚节点更新文件
  
  > 这个操作是可选的
  
    ```shell
    # 每个组织都对应一个锚节点的更新文件
    # go组织锚节点文件
    $ configtxgen -profile ItcastOrgsChannel -outputAnchorPeersUpdate GoMSPanchors.tx -channelID itcastchannel -asOrg OrgGoMSP
    # cpp组织锚节点文件
    $ configtxgen -profile ItcastOrgsChannel -outputAnchorPeersUpdate CppMSPanchors.tx -channelID itcastchannel -asOrg OrgCppMSP
  ```
  
    ```shell
    # 查看生成的文件
    $ tree -L 1
    .
    ├── channel-artifacts
    ├── channel.tx	----------> 生成的通道文件
    ├── configtx.yaml
    ├── CppMSPanchors.tx -----> 生成的cpp组织锚节点文件
    ├── crypto-config
    ├── crypto-config.yaml
    ├── genesis.block --------> 生成的创始块文件
    └── GoMSPanchors.tx	------> 生成的go组织锚节点文件
    ```

## 3. docker-compose文件的编写

### 3.1 客户端角色需要使用的环境变量

```shell
# 客户端docker容器启动之后, go的工作目录
- GOPATH=/opt/gopath	# 不需要修改
# docker容器启动之后, 对应的守护进程的本地套接字, 不需要修改
- CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
- CORE_LOGGING_LEVEL=INFO	# 日志级别 参考hyperledger-fabric.md 3.peer 的环境变量，根据实际情况决定，级别越低，写的日志（io操作，慢）越多，效率越低
- CORE_PEER_ID=cli			# 当前客户端节点的ID, 也是自己这个节点的名字
- CORE_PEER_ADDRESS=peer0.org1.example.com:7051 # 客户端连接的peer节点，哪一个节点都可以
- CORE_PEER_LOCALMSPID=Org1MSP	# 组织ID（客户端连接的peer节点所属组织的组织id）
- CORE_PEER_TLS_ENABLED=true	# 通信是否使用tls加密（客户端与peer节点通信） 这里是false的话，下面TLS相关的内容可以不填；true的话下面与之相关的内容根据你连接的peer节点来填写
- CORE_PEER_TLS_CERT_FILE=		# 证书文件
 /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
- CORE_PEER_TLS_KEY_FILE=		# 私钥文件
 /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
-CORE_PEER_TLS_ROOTCERT_FILE=	# 根证书文件
 /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
# 指定当前客户端的身份，组织的用户里的管理员身份用户
- CORE_PEER_MSPCONFIGPATH=      /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
```

### 3.2 orderer节点需要使用的环境变量

```shell
- ORDERER_GENERAL_LOGLEVEL=INFO	# 日志级别
- ORDERER_GENERAL_LISTENADDRESS=0.0.0.0	# orderer节点监听的地址
- ORDERER_GENERAL_GENESISMETHOD=file	# 创始块的来源, 指定file来源就是文件中
# 创始块对应的文件, 这个不需要改
- ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
- ORDERER_GENERAL_LOCALMSPID=OrdererMSP	# orderer节点所属的组的ID
- ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp	# 当前节点的msp账号路径
# enabled TLS
- ORDERER_GENERAL_TLS_ENABLED=true	# 是否使用tls加密
- ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key	# 私钥
- ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt	# 证书
- ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]			# 根证书
```

### 3.3 peer节点需要使用的环境变量

```shell
- CORE_PEER_ID=peer0.orggo.test.com	# 当前peer节点的名字, 自己起
# 当前peer节点的地址信息
- CORE_PEER_ADDRESS=peer0.orggo.test.com:7051
# 启动的时候, 指定连接谁, 一般写自己就行（启动节点后向哪些节点发起gossip连接, 以加入网络）
- CORE_PEER_GOSSIP_BOOTSTRAP=peer0.orggo.test.com:7051
# 为了被其他节点感知到, 如果不设置别的节点不知有该节点的存在
- CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.orggo.test.com:7051
- CORE_PEER_LOCALMSPID=OrgGoMSP
# docker的本地套接字地址, 不需要改
- CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
# 当前节点属于哪个网络
- CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=network_default
- CORE_LOGGING_LEVEL=INFO
- CORE_PEER_TLS_ENABLED=true
- CORE_PEER_GOSSIP_USELEADERELECTION=true	# 释放自动选举leader节点
- CORE_PEER_GOSSIP_ORGLEADER=false			# 当前不是leader
- CORE_PEER_PROFILE_ENABLED=true	# 在peer节点中有一个profile服务
- CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
- CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
- CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
```

### 3.4 相关配置文件

- **启动docker-compose使用的配置文件** - `docker-compose.yaml`

  ```yaml
  # docker-compose.yaml
  #加上客户端，总共启动了6个容器
  
  version: '2'
  
  #下面有五台远程主机分别映射到本地的：/var/lib/docker/volumes/orderer.itcast.com。。。
  #要结合docker-compose-base.yaml文件才算完成映射
  volumes:
    orderer.itcast.com:
    peer0.orggo.itcast.com:
    peer1.orggo.itcast.com:
    peer0.orgcpp.itcast.com:
    peer1.orgcpp.itcast.com:
  
  #五个docker运行在同一个网络中才能通信
  networks:
    byfn:
  
  services:
  
    orderer.itcast.com: #服务名，跟域名写成了一样，可以自己制定
      extends:
        file:   base/docker-compose-base.yaml
        service: orderer.itcast.com
      container_name: orderer.itcast.com  #容器名，跟域名写成了一样，可以自己制定
      networks:
        - byfn
  
    peer0.orggo.itcast.com:
      container_name: peer0.orggo.itcast.com
      extends:
        file:  base/docker-compose-base.yaml
        service: peer0.orggo.itcast.com
      networks:
        - byfn
  
    peer1.orggo.itcast.com:
      container_name: peer1.orggo.itcast.com
      extends:
        file:  base/docker-compose-base.yaml
        service: peer1.orggo.itcast.com
      networks:
        - byfn
  
    peer0.orgcpp.itcast.com:
      container_name: peer0.orgcpp.itcast.com
      extends:
        file:  base/docker-compose-base.yaml
        service: peer0.orgcpp.itcast.com
      networks:
        - byfn
  
    peer1.orgcpp.itcast.com:
      container_name: peer1.orgcpp.itcast.com
      extends:
        file:  base/docker-compose-base.yaml
        service: peer1.orgcpp.itcast.com
      networks:
        - byfn
  
    cli: #这里是linux终端
      container_name: cli
      image: hyperledger/fabric-tools:latest
      tty: true #终端
      stdin_open: true
      environment:
        - GOPATH=/opt/gopath
        - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
        - CORE_LOGGING_LEVEL=DEBUG
        #- CORE_LOGGING_LEVEL=INFO
        - CORE_PEER_ID=cli
        - CORE_PEER_ADDRESS=peer0.orggo.itcast.com:7051
        - CORE_PEER_LOCALMSPID=OrgGoMSP
        - CORE_PEER_TLS_ENABLED=true
        - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/server.crt
        - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/server.key
        - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/ca.crt
        - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/users/Admin@orggo.itcast.com/msp
      working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
      command: /bin/bash
      volumes:
          - /var/run/:/host/var/run/
          - ./chaincode/:/opt/gopath/src/github.com/chaincode
          - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
          #./channel-artifacts内有什么：看2最下面
          - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
      depends_on: #启动顺序，从上往下启动，最后启动client
        - orderer.itcast.com
        - peer0.orggo.itcast.com
        - peer1.orggo.itcast.com
        - peer0.orgcpp.itcast.com
        - peer1.orgcpp.itcast.com
      networks:
        - byfn
  ```

- 被`docker-compose.yaml`依赖的文件 - `base/docker-compose-base.yaml`

  ```yaml
  # base/docker-compose-base.yaml
  
  version: '2'
  
  services:
  
    orderer.itcast.com:
      container_name: orderer.itcast.com
      image: hyperledger/fabric-orderer:latest
      environment:
        - ORDERER_GENERAL_LOGLEVEL=INFO
        - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
        - ORDERER_GENERAL_GENESISMETHOD=file
        - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
        - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
        - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
        # enabled TLS
        - ORDERER_GENERAL_TLS_ENABLED=true
        - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
        - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
        - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      working_dir: /opt/gopath/src/github.com/hyperledger/fabric
      command: orderer
      volumes:
      - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ../crypto-config/ordererOrganizations/itcast.com/orderers/orderer.itcast.com/msp:/var/hyperledger/orderer/msp
      - ../crypto-config/ordererOrganizations/itcast.com/orderers/orderer.itcast.com/tls/:/var/hyperledger/orderer/tls  #要进行tls通信，所以要有tls证书
      - orderer.itcast.com:/var/hyperledger/production/orderer
      #结合docker-compose.yaml最上面的内容才算完成了整个映射
      # /var/lib/docker/volumes/order.itcast.com
      ports:
        - 7050:7050
  
    peer0.orggo.itcast.com:
      container_name: peer0.orggo.itcast.com
      extends:
        file: peer-base.yaml
        service: peer-base
      environment:
        - CORE_PEER_ID=peer0.orggo.itcast.com
        - CORE_PEER_ADDRESS=peer0.orggo.itcast.com:7051
        - CORE_PEER_GOSSIP_BOOTSTRAP=peer1.orggo.itcast.com:7051
        - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.orggo.itcast.com:7051
        - CORE_PEER_LOCALMSPID=OrgGoMSP
      volumes:
          - /var/run/:/host/var/run/
          - ../crypto-config/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/msp:/etc/hyperledger/fabric/msp
          - ../crypto-config/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls:/etc/hyperledger/fabric/tls
          - peer0.orggo.itcast.com:/var/hyperledger/production
      ports:
        - 7051:7051 #正常通信的端口
        - 7053:7053 #fabric里面有一些事件，如果事件触发了，会通过这个端口传输数据
  
    peer1.orggo.itcast.com:
      container_name: peer1.orggo.itcast.com
      extends:
        file: peer-base.yaml
        service: peer-base
      environment:
        - CORE_PEER_ID=peer1.orggo.itcast.com
        - CORE_PEER_ADDRESS=peer1.orggo.itcast.com:7051
        - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.orggo.itcast.com:7051
        - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.orggo.itcast.com:7051
        - CORE_PEER_LOCALMSPID=OrgGoMSP
      volumes:
          - /var/run/:/host/var/run/
          - ../crypto-config/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/msp:/etc/hyperledger/fabric/msp
          - ../crypto-config/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/tls:/etc/hyperledger/fabric/tls
          - peer1.orggo.itcast.com:/var/hyperledger/production
  
      ports:
        - 8051:7051
        - 8053:7053
  
    peer0.orgcpp.itcast.com:
      container_name: peer0.orgcpp.itcast.com
      extends:
        file: peer-base.yaml
        service: peer-base
      environment:
        - CORE_PEER_ID=peer0.orgcpp.itcast.com
        - CORE_PEER_ADDRESS=peer0.orgcpp.itcast.com:7051
        - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.orgcpp.itcast.com:7051
        - CORE_PEER_GOSSIP_BOOTSTRAP=peer1.orgcpp.itcast.com:7051
        - CORE_PEER_LOCALMSPID=OrgCppMSP
      volumes:
          - /var/run/:/host/var/run/
          - ../crypto-config/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/msp:/etc/hyperledger/fabric/msp
          - ../crypto-config/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/tls:/etc/hyperledger/fabric/tls
          - peer0.orgcpp.itcast.com:/var/hyperledger/production
      ports:
        - 9051:7051
        - 9053:7053
  
    peer1.orgcpp.itcast.com:
      container_name: peer1.orgcpp.itcast.com
      extends:
        file: peer-base.yaml
        service: peer-base
      environment:
        - CORE_PEER_ID=peer1.orgcpp.itcast.com
        - CORE_PEER_ADDRESS=peer1.orgcpp.itcast.com:7051
        - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.orgcpp.itcast.com:7051
        - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.orgcpp.itcast.com:7051
        - CORE_PEER_LOCALMSPID=OrgCppMSP
      volumes:
          - /var/run/:/host/var/run/
          - ../crypto-config/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/msp:/etc/hyperledger/fabric/msp
          - ../crypto-config/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/tls:/etc/hyperledger/fabric/tls
          - peer1.orgcpp.itcast.com:/var/hyperledger/production
      ports:
        - 10051:7051
        - 10053:7053
  ```

- 被 ``docker-compose-base.yaml` 依赖的文件 - `base/peer-base.yaml`

  ```yaml
  # base/peer-base.yaml
  
  version: '2'
  
  services:
    peer-base:
      image: hyperledger/fabric-peer:latest
      environment:
        - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
        # the following setting starts chaincode containers on the same
        # bridge network as the peers
        # https://docs.docker.com/compose/networking/
        - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=xxxx_byfn
        - CORE_LOGGING_LEVEL=INFO
        #- CORE_LOGGING_LEVEL=DEBUG
        - CORE_PEER_TLS_ENABLED=true
        - CORE_PEER_GOSSIP_USELEADERELECTION=true
        - CORE_PEER_GOSSIP_ORGLEADER=false
        - CORE_PEER_PROFILE_ENABLED=true
        - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
        - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
        - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
      command: peer node start
  
  ```

### 3.5 启动docker-compose

![1542248425558](D:/%E9%BB%91%E9%A9%AC%E5%8C%BA%E5%9D%97%E9%93%BE/%E8%B5%84%E6%96%99/%E7%AC%AC%E4%B8%83%E9%98%B6%E6%AE%B5/hyperledger/day03-%E9%80%9A%E9%81%93%E6%93%8D%E4%BD%9C%E5%92%8C%E6%99%BA%E8%83%BD%E5%90%88%E7%BA%A6/01-%E6%95%99%E5%AD%A6%E8%B5%84%E6%96%99/assets/1542248425558.png)

```shell
CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=${COMPOSE_PROJECT_NAME}_byfn -> _byfn
创建的网络叫: itcast_byfn
	- byfn: 网络名
	- itcast: docker-compose.yaml所在的目录
```

检测网络是否正常启动了:

```shell
# 在docker-compose.yaml 文件目录下执行下边命令
$ docker-compose ps
         Name                 Command       State                        Ports                      
----------------------------------------------------------------------------------------------------
cli                       /bin/bash         Up                                                      
orderer.itcast.com        orderer           Up      0.0.0.0:7050->7050/tcp                          
peer0.orgcpp.itcast.com   peer node start   Up      0.0.0.0:9051->7051/tcp, 0.0.0.0:9053->7053/tcp  
peer0.orggo.itcast.com    peer node start   Up      0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp  
peer1.orgcpp.itcast.com   peer node start   Up      0.0.0.0:10051->7051/tcp, 0.0.0.0:10053->7053/tcp
peer1.orggo.itcast.com    peer node start   Up      0.0.0.0:8051->7051/tcp, 0.0.0.0:8053->7053/tcp
```

## 4.  Peer操作命令（先看第五部分开头，再看4.peer操作命令）

peer命令只能在peer镜像启动后查看，或者是cli启动之后查看

### 4.1 创建通道

```shell
$ peer channel create [flags], 常用参数为:
	`-o, --orderer: orderer节点的地址
	`-c, --channelID: 要创建的通道的ID, 必须小写, 在250个字符以内
	`-f, --file: 由configtxgen 生成的通道文件, 用于提交给orderer
	-t, --timeout: 创建通道的超时时长, 默认为5s
	`--tls: 通信时是否使用tls加密
	`--cafile: 当前orderer节点pem格式的tls证书文件, 要使用绝对路径.
# orderer节点pem格式的tls证书文件路径参考: 
crypto-config/ordererOrganizations/itcast.com/orderers/orderer.itcast.com/msp/tlscacerts/tlsca.itcast.com-cert.pem
# example
$ peer channel create -o orderer节点地址:端口 -c 通道名 -f 通道文件 --tls true --cafile orderer节点pem格式的证书文件
	- orderer节点地址: 可以是IP地址或者域名
	- orderer节点监听的是7050端口
$ peer channel create -o orderer.itcast.com:7050 -c itcastchannel -f ./channel-artifacts/channel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/itcast.com/msp/tlscacerts/tlsca.itcast.com-cert.pem
# 在当前工作目录下生成一个文件: 通道名.block, 本例: itcastchannel.block
$ ls
channel-artifacts  crypto  `itcastchannel.block` --> 生成的文件
```

### **4.2 加入通道**

```shell
$ peer channel join[flags], 常用参数为:
	`-b, --blockpath: 通过 peer channel create 命令生成的通道文件 
# example
$ peer channel join -b 生成的通道block文件
$ peer channel join -b itcastchannel.block 
```

### 补充：

其他的节点，加入通道

```
# 第1个节点 Go组织的 peer0
export CORE_PEER_ADDRESS=peer0.orggo.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgGoMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/users/Admin@orggo.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/server.key

# 第2个节点 Go组织的 peer1
export CORE_PEER_ADDRESS=peer1.orggo.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgGoMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/users/Admin@orggo.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/tls/server.key

# 第3个节点 Cpp组织的 peer0 注意：cli为自己定义的客户端名称，根据自己的命名进行修改
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer0.orgcpp.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgCppMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/users/Admin@orgcpp.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/tls/server.key

# 第4个节点 Cpp组织的 peer1 注意：cli为自己定义的客户端名称，根据自己的命名进行修改
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer1.orgcpp.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgCppMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/users/Admin@orgcpp.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/tls/server.key
```



### 4.3 更新锚节点

```shell
$ peer channel update [flags], 常用参数为:
	`-o, --orderer: orderer节点的地址
	`-c, --channelID: 要创建的通道的ID, 必须小写, 在250个字符以内
	`-f, --file: 由configtxgen 生成的组织锚节点文件, 用于提交给orderer
	`--tls: 通信时是否使用tls加密
	`--cafile: 当前orderer节点pem格式的tls证书文件, 要使用绝对路径.
# orderer节点pem格式的tls证书文件路径参考: 
crypto-config/ordererOrganizations/itcast.com/orderers/orderer.itcast.com/msp/tlscacerts/tlsca.itcast.com-cert.pem
# example
$ peer channel update -o orderer节点地址:端口 -c 通道名 -f 锚节点更新文件 --tls true --cafile orderer节点pem格式的证书文件
```

### **4.4 安装链码**

```shell
$ peer chaincode install [flags], 常用参数为:
	-c, --ctor: JSON格式的构造参数, 默认是"{}"
	`-l, --lang: 编写chaincode的编程语言, 默认值是 golang
	`-n, --name: chaincode的名字
	`-p, --path: chaincode源代码的目录, 从 $GOPATH/src 路径后开始写
	`-v, --version: 当前操作的chaincode的版本, 适用这些命令install/instantiate/upgrade
# example
$ peer chaincode install -n 链码的名字 -v 链码的版本 -l 链码的语言 -p 链码的位置
	- 链码名字自己起
	- 链码的版本, 自己根据实际情况指定
$ peer chaincode install -n testcc -v 1.0 -l golang -p github.com/chaincode
```

### **4.5 链码初始化**

```shell
$ peer chaincode instantiate [flags], 常用参数为:
	`-C，--channelID：当前命令运行的通道，默认值是“testchainid"。
	`-c, --ctor：JSON格式的构造参数，默认值是“{}"
	`-l，--lang：编写Chaincode的编程语言，默认值是golang
	`-n，--name：Chaincode的名字。
	`-P，--policy：当前Chaincode的背书策略。
	`-v，--version：当前操作的Chaincode的版本，适用于install/instantiate/upgrade等命令
	`--tls: 通信时是否使用tls加密
	`--cafile: 当前orderer节点pem格式的tls证书文件, 要使用绝对路径.
	 
# example
# -c '{"Args":["init","a","100","b","200"]}' 
# -P "AND ('OrgGoMSP.member', 'OrgCppMSP.member')"

$ peer chaincode instantiate -o orderer节点地址:端口 --tls true --cafile orderer节点pem格式的证书文件 -C 通道名称 -n 链码名称 -l 链码语言 -v 链码版本 -c 链码Init函数调用 -P 背书策略

$ peer chaincode instantiate -o orderer.itcast.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/itcast.com/msp/tlscacerts/tlsca.itcast.com-cert.pem  -C itcastchannel -n testcc -l golang -v 1.0 -c '{"Args":["init","a","100","b","200"]}' -P "AND ('OrgGoMSP.member', 'OrgCppMSP.member')"
```

### **4.6 查询**

```shell
$ peer chaincode query [flags], 常用参数为:
	`-n，--name：Chaincode的名字。
	`-C，--channelID：当前命令运行的通道，默认值是“testchainid"
	`-c, --ctor：JSON格式的构造参数，默认值是“{}"
	-x，--hex：是否对输出的内容进行编码处理
	-r，--raw：是否输出二进制内容
	-t, --tid: 指定当前查询的编号
# example
# '{"Args":["query","a"]}'
$ peer chaincode query -C 通道名称 -n 链码名称 -c 链码调用
```

### **4.7 交易**

```shell
$ peer chaincode invoke [flags], 常用参数为:
	`-o, --orderer: orderer节点的地址
	`-C，--channelID：当前命令运行的通道，默认值是“testchainid"
	`-c, --ctor：JSON格式的构造参数，默认值是“{}"
	`-n，--name：Chaincode的名字
	`--tls: 通信时是否使用tls加密
	`--cafile: 当前orderer节点pem格式的tls证书文件, 要使用绝对路径.
	`--peerAddresses: 指定要连接的peer节点的地址
	`--tlsRootCertFiles: 连接的peer节点的TLS根证书
# 连接的peer节点的TLS根证书查找路径参考:
/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/ca.crt
# example
# -c '{"Args":["invoke","a","b","10"]}'
$ peer chaincode invoke -o orderer节点地址:端口 --tls true --cafile orderer节点pem格式的证书文件 -C 通道名称 -n 链码名称 --peerAddresses 背书节点1:端口 --tlsRootCertFiles 背书节点1的TLS根证书    --peerAddresses 背书节点2:端口 --tlsRootCertFiles 背书节点2的TLS根证书 -c 交易链码调用
```

## 5. 通过客户端操作各节点（先看第五部分开头，再看4.peer操作命令）

客户端对Peer节点的操作流程:

- **创建通道, 通过客户端节点来完成**

  ```shell
  # 在宿主机
  $ docker-compose ps
           Name                 Command       State                        Ports                      
  ----------------------------------------------------------------------------------------------------
  cli                       /bin/bash         Up                                                      
  orderer.itcast.com        orderer           Up      0.0.0.0:7050->7050/tcp                          
  peer0.orgcpp.itcast.com   peer node start   Up      0.0.0.0:9051->7051/tcp, 0.0.0.0:9053->7053/tcp  
  peer0.orggo.itcast.com    peer node start   Up      0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp  
  peer1.orgcpp.itcast.com   peer node start   Up      0.0.0.0:10051->7051/tcp, 0.0.0.0:10053->7053/tcp
  peer1.orggo.itcast.com    peer node start   Up      0.0.0.0:8051->7051/tcp, 0.0.0.0:8053->7053/tcp 
  # 进入到客户端对用的容器中
  $ docker exec -it cli /bin/bash
  ```

- 1.将每个组织的每个节点都加入到通道中  -> 客户端来完成的

  - 一个客户端同时只能连接一个peer节点（配置文件中的环境变量已经指定了要连接的peer节点），如果想要修改客户端连接的peer节点，进入客户端修改环境变量就可以了

- 2.给每个peer节点安装智能合约 -> 链代码(程序: go, node.js, java)

- 3.对智能合约进行初始化 , 对应智能合约中的 Init 函数

  - <font color="red">只需要在任意节点初始化一次, 数据会自动同步的各个组织的各个节点</font>

- 4.对数据进行查询 -> 读

- 5.对数据进行调用 -> 写

> 经过前面的讲解我们都知道, 一个客户端只能连接一个指定的节点, 如果想要该客户端连接其他节点, 那么就必须修改当前客户端中相关的环境变量

### 5.1 相关环境变量

```shell
# 第1个节点 Go组织的 peer0
export CORE_PEER_ADDRESS=peer0.orggo.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgGoMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/users/Admin@orggo.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer0.orggo.itcast.com/tls/server.key

# 第2个节点 Go组织的 peer1
export CORE_PEER_ADDRESS=peer1.orggo.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgGoMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/users/Admin@orggo.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.itcast.com/peers/peer1.orggo.itcast.com/tls/server.key

# 第3个节点 Cpp组织的 peer0
export CORE_PEER_ADDRESS=peer0.orgcpp.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgCppMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/users/Admin@orgcpp.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer0.orgcpp.itcast.com/tls/server.key

# 第4个节点 Cpp组织的 peer1
export CORE_PEER_ADDRESS=peer1.orgcpp.itcast.com:7051
export CORE_PEER_LOCALMSPID=OrgCppMSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/users/Admin@orgcpp.itcast.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.itcast.com/peers/peer1.orgcpp.itcast.com/tls/server.key
```

> 创建通道的操作只要登录到客户端容器中就可以进行操作.

### 5.1 对peer0.OrgGo的操作

- 要保证客户端操作的是peer0.OrgGo
  - 可以查看:  `echo $CORE_PEER_ADDRESS`
- 将当前节点加入到通道中
  - `peer channel join -b xxx.block`
- 安装链代码
  - `peer chaincode install [flags]`
- 链代码的初始化  -> 只需要做一次
  - `peer chaincode instantiate [flag]`
- 查询/调用

### 5.2 对peer1.OrgGo的操作

- 要保证客户端操作的是peer1.OrgGo
  - 可以查看:  `echo $CORE_PEER_ADDRESS`
  - 不是修改环境变量
- 将当前节点加入到通道中
  - `peer channel join -b xxx.block`
- 安装链代码
  - `peer chaincode install [flags]`
- 查询/调用

### 5.3 对peer0.OrgCpp的操作

- 要保证客户端操作的是peer1.OrgGo
  - 可以查看:  `echo $CORE_PEER_ADDRESS`
  - 不是修改环境变量
- 将当前节点加入到通道中
  - `peer channel join -b xxx.block`
- 安装链代码
  - `peer chaincode install [flags]`
- 查询/调用

### 5.4 对peer1.OrgCpp的操作

- 要保证客户端操作的是peer1.OrgGo
  - 可以查看:  `echo $CORE_PEER_ADDRESS`
  - 不是修改环境变量
- 将当前节点加入到通道中
  - `peer channel join -b xxx.block`
- 安装链代码
  - `peer chaincode install [flags]`
- 查询/调用

## 7. 智能合约

### 7.1 常识

- 链代码的**包名**的指定

  ```go
  // xxx.go
  package main
  ```

- 必须要引入的包

  ```go
  // go get github.com/hyperledger/fabric/core/chaincode/shim
  import (
      // 客户端需要和 Fabric框架通信
      "github.com/hyperledger/fabric/core/chaincode/shim"
      pb "github.com/hyperledger/fabric/protos/peer"
  )
  ```

- 链码的书写要求

  ```go
  // 自定义一个结构体 - 类, 基于这个类实现一些接口函数
  type Test struct {
      // 空着即可
  }
  func (t* Test) Init(stub ChaincodeStubInterface) pb.Response;
  func (t* Test) Invoke(stub ChaincodeStubInterface) pb.Response;
  ```

- <font color="red">链码 API 查询</font>

  ```http
  https://godoc.org/github.com/hyperledger/fabric/core/chaincode/shim
  ```

### 7.2 常用接口

```go
// 操作账本的
// shim -> ChaincodeStubInterface
```



## 示例链码

```go
package main

import (
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

```







