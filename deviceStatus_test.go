package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xmidt-org/webpa-common/v2/convey"
	"github.com/xmidt-org/webpa-common/v2/device"
)

func deviceWithMetadata(metadata *device.Metadata) *device.MockDevice {
	d := new(device.MockDevice)
	d.On("ID").Return(device.ID("mac:112233445566"))
	d.On("Metadata").Return(metadata)
	d.On("Convey").Return(convey.C{})
	d.On("ConveyCompliance").Return(convey.Full)
	return d
}

func TestOnlineMessageCarriesTheAddressOfTheConnection(t *testing.T) {
	assert := assert.New(t)

	metadata := new(device.Metadata)
	metadata.Store(ConnectIPMetadataKey, "203.0.113.7")
	_, message := newOnlineMessage("dns:talaria", deviceWithMetadata(metadata))

	assert.Equal("203.0.113.7", message.Metadata["/connect-ip"])
}

func TestOnlineMessageWithoutAnAddressHasNone(t *testing.T) {
	assert := assert.New(t)

	_, message := newOnlineMessage("dns:talaria", deviceWithMetadata(new(device.Metadata)))

	_, found := message.Metadata["/connect-ip"]
	assert.False(found)
}

func TestTheAddressIsNotInTheStatusMetadataOfOtherEvents(t *testing.T) {
	assert := assert.New(t)

	metadata := new(device.Metadata)
	metadata.Store(ConnectIPMetadataKey, "203.0.113.7")

	// statusMetadata is what the offline event uses too: the address goes with the online event only
	_, found := statusMetadata(deviceWithMetadata(metadata))["/connect-ip"]
	assert.False(found)
}
