package epoll

import (
	"github.com/lixianmin/road/core"
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
	input   *core.Buffer
	backup  *core.Buffer
}

func NewWebReaderWriter(conn net.Conn, watcher *gaio.Watcher) *WebReaderWriter {
	var my = &WebReaderWriter{
		conn:    conn,
		watcher: watcher,
		input:   &core.Buffer{},
		backup:  &core.Buffer{},
	}

	return my
}

// onReceiveData()与下面的Read()是在同一个线程中调用的，不存在并发问题
func (my *WebReaderWriter) onReceiveData(buff []byte) {
	_, _ = my.input.Write(buff)
}

func (my *WebReaderWriter) InputSize() int {
	return my.input.Len()
}

func (my *WebReaderWriter) TakeSnapshot() {
	my.backup.Reset()
	my.backup.Write(my.input.Bytes())
}

func (my *WebReaderWriter) Rollback() {
	my.input.Reset()
	my.input.Write(my.backup.Bytes())
}

func (my *WebReaderWriter) Read(p []byte) (n int, err error) {
	n, err = my.input.Read(p)
	return n, err
}

func (my *WebReaderWriter) Write(p []byte) (n int, err error) {
	return len(p), my.watcher.Write(my, my.conn, p)
}
