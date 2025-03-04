package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/platforms/parrot/minidrone"

	"github.com/joek/kyma-drone/pkg/drone"
	connector "github.com/joek/kyma-drone/pkg/kyma-connector"
	"github.com/joek/kyma-drone/pkg/models"

	"github.com/go-openapi/runtime/middleware"
	"github.com/joek/kyma-drone/pkg/restapi/operations"
)

// PublicShippingDroneHandler is handling request to ship the package
type PublicShippingDroneHandler struct {
	drone drone.Drone
	conn  *connector.KymaConnector
}

// Handle http Handler to Up drones
func (h PublicShippingDroneHandler) Handle(params operations.ShipPackageParams) middleware.Responder {
	orderCoder := params.Value.OrderCode
	h.drone.Once(minidrone.Landed, func(data interface{}) {
		log.Println("Landed")
		h.conn.SendEvent(json.RawMessage([]byte("{\"orderCode\": \""+*orderCoder+"\"}")), "drone.shipped", "v1")
	})
	h.drone.Once(minidrone.Hovering, func(data interface{}) {
		time.AfterFunc(10*time.Second, func() {
			log.Println("Landing")
			h.drone.Land()
		})
	})
	err := h.drone.TakeOff()
	if err != nil {
		c := int32(20)
		m := fmt.Sprintf("Error: %s", err)
		er := operations.NewUpDroneDefault(-10)
		er.SetPayload(&models.ErrorModel{
			Code:    &c,
			Message: &m,
		})
		return er
	}

	return &operations.UpDroneNoContent{}
}

// NewPublicShippingDroneHandler is creating a new Up Handler
func NewPublicShippingDroneHandler(d drone.Drone, conn *connector.KymaConnector) PublicShippingDroneHandler {
	return PublicShippingDroneHandler{d, conn}
}
