//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"sync"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	redisClient "github.com/edgexfoundry/edgex-go/internal/pkg/db/redis"
	"github.com/gomodule/redigo/redis"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/errors"
	model "github.com/edgexfoundry/go-mod-core-contracts/v2/models"

	"github.com/google/uuid"
)

var currClient *Client // a singleton so Readings can be de-referenced
var once sync.Once

type Client struct {
	*redisClient.Client
	loggingClient logger.LoggingClient
}

func NewClient(config db.Configuration, logger logger.LoggingClient) (*Client, errors.EdgeX) {
	var err error
	dc := &Client{}
	dc.Client, err = redisClient.NewClient(config, logger)
	dc.loggingClient = logger
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "redis client creation failed", err)
	}

	return dc, nil
}

// CloseSession closes the connections to Redis
func (c *Client) CloseSession() {
	c.Pool.Close()

	currClient = nil
	once = sync.Once{}
}

// Add a new event
func (c *Client) AddEvent(e model.Event) (model.Event, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if e.Id != "" {
		_, err := uuid.Parse(e.Id)
		if err != nil {
			return model.Event{}, errors.NewCommonEdgeX(errors.KindInvalidId, "uuid parsing failed", err)
		}
	}

	return addEvent(conn, e)
}

// Get an event by id
func (c *Client) EventById(id string) (event model.Event, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	event, edgeXerr = eventById(conn, id)
	if edgeXerr != nil {
		if edgeXerr == redis.ErrNil {
			return event, errors.NewCommonEdgeXWrapper(edgeXerr)
		}
		return event, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// Add a new device profle
func (c *Client) AddDeviceProfile(dp model.DeviceProfile) (model.DeviceProfile, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if dp.Id != "" {
		_, err := uuid.Parse(dp.Id)
		if err != nil {
			return model.DeviceProfile{}, errors.NewCommonEdgeX(errors.KindInvalidId, "ID failed UUID parsing", err)
		}
	} else {
		dp.Id = uuid.New().String()
	}

	return addDeviceProfile(conn, dp)
}

// Update a new device profile
func (c *Client) UpdateDeviceProfile(dp model.DeviceProfile) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateDeviceProfile(conn, dp)
}

// Add a new device service
func (c *Client) AddDeviceService(ds model.DeviceService) (model.DeviceService, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(ds.Id) == 0 {
		ds.Id = uuid.New().String()
	}

	return addDeviceService(conn, ds)
}

// Get a device service by name
func (c *Client) GetDeviceServiceByName(name string) (deviceService model.DeviceService, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceService, edgeXerr = deviceServiceByName(conn, name)
	if edgeXerr != nil {
		if edgeXerr == redis.ErrNil {
			return deviceService, errors.NewCommonEdgeXWrapper(edgeXerr)
		}
		return deviceService, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// Get a device profile by name
func (c *Client) GetDeviceProfileByName(name string) (deviceProfile model.DeviceProfile, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceProfile, edgeXerr = deviceProfileByName(conn, name)
	if edgeXerr != nil {
		if edgeXerr == redis.ErrNil {
			return deviceProfile, errors.NewCommonEdgeXWrapper(edgeXerr)
		}
		return deviceProfile, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// Delete a device profile by id
func (c *Client) DeleteDeviceProfileById(id string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceProfileById(conn, id)
	if edgeXerr != nil {
		if edgeXerr == redis.ErrNil {
			return errors.NewCommonEdgeXWrapper(edgeXerr)
		}
		return errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return nil
}

// Delete a device profile by name
func (c *Client) DeleteDeviceProfileByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceProfileByName(conn, name)
	if edgeXerr != nil {
		if edgeXerr == redis.ErrNil {
			return errors.NewCommonEdgeXWrapper(edgeXerr)
		}
		return errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return nil
}