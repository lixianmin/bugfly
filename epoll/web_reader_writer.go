package epoll

import (
	"bytes"
	"github.com/xtaci/gaio"
	"net"
)

/********************************************************************
created:    2020-12-07
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WebReaderWriter struct {
	conn    net.Conn
	watcher *gaio.Watcher
	input   *bytes.Buffer
}

func NewWebReaderWriter(conn net.Conn, watcher *gaio.Watcher) *WebReaderWriter {
	var my = &WebReaderWriter{
		conn:    conn,
		watcher: watcher,
		input:   gBufferPool.Get(),
	}

	return my
}

// onReceiveData()与下面的Read()是在同一个线程中调用的，不存在并发问题
func (my *WebReaderWriter) onReceiveData(buff []byte) {
	_, _ = my.input.Write(buff)
}

func (my *WebReaderWriter) Read(p []byte) (n int, err error) {
	n, err = my.input.Read(p)
	my.input = checkSwapBuffer(my.input)

	return n, err
}

func (my *WebReaderWriter) Write(p []byte) (n int, err error) {
	return len(p), my.watcher.Write(my, my.conn, p)
}