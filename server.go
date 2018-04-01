package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/juju/ratelimit"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

func runServer(cmd *cobra.Command, args []string) error {
	logger := log15.New()

	d := &dispatch{
		Bucket:        ratelimit.NewBucket(1*time.Second, 10),
		Subscriptions: map[string]subscriptions{},
		Logger:        logger,
	}

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Fprint(w, "Not much to see here")
	})

	router.ServeFiles("/static/*filepath", http.Dir("./static"))

	router.GET("/status", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Fprint(w, "OK")
	})

	router.GET("/remote/:id", serveFrontend)
	router.POST("/trigger/:id", d.Trigger) // code in post body: used by ifttt
	router.POST("/send/:id/:code", d.Send) // code in url: used by frontend
	router.GET("/subscribe/:id", d.Subscribe)

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}
	logger.Info("starting server", "port", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func serveFrontend(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := log15.New()

	tmpl, err := template.New("remote.html").ParseGlob("templates/*.html")
	if err != nil {
		logger.Error("failed to parse template", "error", err)
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	data := map[string]string{"ID": params.ByName("id")}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("failed to execute template", "error", err)
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
		return
	}
}
