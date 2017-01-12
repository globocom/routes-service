package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/osrg/gobgp/client"
	"github.com/osrg/gobgp/config"
)

type Neighbor struct {
	PeerAs          uint32
	NeighborAddress string
}

func main() {
	gobgp, err := client.New("192.168.33.10:50051")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.POST("/neighbor", func(c echo.Context) error {
		var neighbor Neighbor

		if err := c.Bind(&neighbor); err != nil {
			return err
		}

		peer := &config.Neighbor{
			Config: config.NeighborConfig{
				PeerAs:          neighbor.PeerAs,
				NeighborAddress: neighbor.NeighborAddress,
			},
		}
		err := gobgp.AddNeighbor(peer)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "ok")
	})

	e.Logger.Fatal(e.Start(":8888"))
}
