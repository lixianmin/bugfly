package road

import (
	"github.com/lixianmin/road/conn/message"
	"github.com/lixianmin/road/conn/packet"
	"github.com/lixianmin/road/logger"
	"github.com/lixianmin/road/util"
)

/********************************************************************
created:    2020-08-31
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func (my *Session) Push(route string, v interface{}) error {
	if my.wc.IsClosed() {
		return nil
	}

	var payload, err = util.SerializeOrRaw(my.serializer, v)
	var msg = message.Message{Type: message.Push, Route: route, Data: payload}
	return my.sendMessageMayError(msg, err)
}

// 强踢下线
func (my *Session) Kick() error {
	if my.wc.IsClosed() {
		return nil
	}

	p, err := my.packetEncoder.Encode(packet.Kick, nil)
	if err != nil {
		return err
	}

	my.sendBytes(p)
	return nil
}

func (my *Session) sendMessageMayError(msg message.Message, err error) error {
	if err != nil {
		msg.Err = true
		logger.Info("process failed, route=%s, err=%q", msg.Route, err.Error())

		// err需要支持json序列化的话，就不能是一个简单的字符串
		var errWrap = checkCreateError(err)

		var err1 error
		msg.Data, err1 = util.SerializeOrRaw(my.serializer, errWrap)
		if err1 != nil {
			logger.Info("serialize failed, route=%s, err1=%q", msg.Route, err1.Error())
			return err1
		}
	}

	data, err2 := my.packetEncodeMessage(&msg)
	if err2 != nil {
		logger.Info("send failed, route=%s, err2=%q", msg.Route, err2.Error())
		return err2
	}

	my.sendBytes(data)
	return nil
}

func (my *Session) sendBytes(data []byte) {
	if len(data) == 0 {
		return
	}

	select {
	case my.sendingChan <- data:
	case <-my.wc.C():
	}
}

func (my *Session) packetEncodeMessage(msg *message.Message) ([]byte, error) {
	data, err := my.messageEncoder.Encode(msg)
	if err != nil {
		return nil, err
	}

	// packet encode
	p, err := my.packetEncoder.Encode(packet.Data, data)
	if err != nil {
		return nil, err
	}

	return p, nil
}
