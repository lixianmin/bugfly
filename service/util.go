package service

import (
	"fmt"
	"github.com/lixianmin/road/component"
	"github.com/lixianmin/road/route"
)

/********************************************************************
created:    2020-08-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func GetHandler(rt *route.Route) (*component.Handler, error) {
	handler, ok := handlers[rt.Short()]
	if !ok {
		e := fmt.Errorf("pitaya/handler: %s not found", rt.String())
		return nil, e
	}
	return handler, nil

}
