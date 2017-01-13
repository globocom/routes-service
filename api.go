package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type API struct {
	echo       *echo.Echo
	bgpService *BgpService
}

func NewAPI(bgpService *BgpService) (*API, error) {
	echo := echo.New()
	api := API{
		echo:       echo,
		bgpService: bgpService,
	}

	echo.POST("/neighbor", api.addNeighborHandler)

	return &api, nil
}

func (api *API) Start() {
	api.echo.Logger.Fatal(api.echo.Start(":8888"))
}

func (api *API) addNeighborHandler(c echo.Context) error {
	var neighbor Neighbor

	if err := c.Bind(&neighbor); err != nil {
		log.Println(err)
		return err
	}

	if err := api.bgpService.AddNeighbor(neighbor); err != nil {
		log.Println(err)
		return err
	}

	return c.String(http.StatusOK, "ok")
}
