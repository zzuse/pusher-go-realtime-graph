package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pusher "github.com/pusher/pusher-http-go"
)

var client = pusher.Client{
	AppId:   "561719",
	Key:     "6906f2e1d3d5a64af997",
	Secret:  "c9df10d84b509d963033",
	Cluster: "ap1",
	Secure:  true,
}

type visitsData struct {
	Pages int
	Count int
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "public/index.html")
	e.File("/style.css", "public/sytle.css")
	e.File("/app.js", "public/app.js")
	e.GET("/simulate", simulate)

	e.Logger.Fatal(e.Start(":9000"))
}

func setInterval(ourFunc func(), milliseconds int, async bool) chan bool {
	interval := time.Duration(milliseconds) * time.Millisecond
	ticker := time.NewTicker(interval)
	clear := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go ourFunc()
				} else {
					ourFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()
	return clear
}

func simulate(c echo.Context) error {
	setInterval(func() {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		newVisitsData := visitsData{
			Pages: r1.Intn(100),
			Count: r1.Intn(100),
		}
		client.Trigger("visitorsCount", "addNumber", newVisitsData)
	}, 2500, true)
	return c.String(http.StatusOK, "Simulation begun")
}
