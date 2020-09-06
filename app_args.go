package road

import (
	"github.com/lixianmin/logo"
	"github.com/lixianmin/road/epoll"
	"time"
)

/********************************************************************
created:    2020-08-31
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type AppArgs struct {
	Acceptor         *epoll.Acceptor
	HeartbeatTimeout time.Duration // 心跳超时时间
	DataCompression  bool          // 数据是否压缩
	Logger           logo.ILogger  // 自定义日志对象，默认只输出到控制台
}
