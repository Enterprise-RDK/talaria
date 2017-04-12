package main

import (
	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/logging"
	"github.com/Comcast/webpa-common/wrp"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
)

func NewPrimaryHandler(logger logging.Logger, manager device.Manager, v *viper.Viper) (http.Handler, error) {
	poolFactory, err := wrp.NewPoolFactory(v.Sub(wrp.ViperKey))
	if err != nil {
		return nil, err
	}

	handler := mux.NewRouter()

	handler.Handle("/device", &device.MessageHandler{
		Logger:   logger,
		Decoders: poolFactory.NewDecoderPool(wrp.JSON),
		Router:   manager,
	}).
		Methods("POST", "PATCH").
		Headers("Content-Type", wrp.JSON.ContentType())

	handler.Handle("/device", &device.MessageHandler{
		Logger:   logger,
		Decoders: poolFactory.NewDecoderPool(wrp.Msgpack),
		Router:   manager,
	}).
		Methods("POST", "PATCH").
		Headers("Content-Type", wrp.Msgpack.ContentType())

	handler.Handle("/connect", &device.ConnectHandler{
		Logger:    logger,
		Connector: manager,
	})

	return handler, nil
}
