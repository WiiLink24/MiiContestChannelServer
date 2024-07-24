package main

import (
	"MiiContestChannelServer/middleware"
	"MiiContestChannelServer/mii"
	"MiiContestChannelServer/webpanel"
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var (
	ctx  = context.Background()
	pool *pgxpool.Pool
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Everybody Votes Channel server has encountered a fatal error! Reason: %v\n", err)
	}
}

func main() {
	// Get config
	config := GetConfig()

	provider, err := oidc.NewProvider(ctx, config.OIDCConfig.Provider)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}

	authConfig := &webpanel.AppAuthConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     config.OIDCConfig.ClientID,
			ClientSecret: config.OIDCConfig.ClientSecret,
			RedirectURL:  config.OIDCConfig.RedirectURL,
			Scopes:       config.OIDCConfig.Scopes,
			Endpoint:     provider.Endpoint(),
		},
		Provider: provider,
	}

	// Start SQL
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	checkError(err)

	// Set up HTTP
	r := gin.Default()
	if gin.Mode() == gin.DebugMode {
		r.Static("/assets", "./assets") // Serve static files
	}
	r.LoadHTMLGlob("templates/*")
	gin.SetMode(gin.ReleaseMode)
	panel := webpanel.WebPanel{
		Pool:       pool,
		Ctx:        ctx,
		Config:     config,
		AuthConfig: authConfig,
	}

	r.GET("/panel/login", panel.LoginPage)
	r.GET("/panel/start", panel.StartPanelHandler)
	r.GET("/panel/authorize", panel.FinishPanelHandler)

	auth := r.Group("/panel")
	if config.AuthMode {
		auth.Use(middleware.AuthenticationMiddleware())
	}
	{
		auth.GET("/admin", panel.AdminPage)
		auth.GET("/contests", panel.ViewContests)
		auth.POST("/contests", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/panel/contests")
		})
		auth.GET("/contests/add", panel.AddContest)
		auth.POST("/contests/add", panel.AddContestPOST)
		auth.POST("/contests/delete/:contest_id", panel.DeleteContest)
		auth.GET("/contests/edit/:contest_id", panel.EditContest)
		auth.POST("/contests/edit/:contest_id", panel.EditContestPOST)
		auth.GET("/plaza", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/panel/plaza/1")
		})
		auth.GET("/plaza/:page", panel.ViewPlaza)
		auth.GET("/plaza/top", panel.ViewPlazaTop50)
		auth.GET("/plaza/new", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/panel/plaza/new/1")
		})
		auth.GET("/plaza/new/:page", panel.ViewPlazaNew)
		auth.POST("/plaza/search", panel.SearchPlaza)
		auth.GET("/plaza/details/:entry_id", panel.ViewMiiDetails)
		auth.POST("/plaza/delete/:entry_id", panel.DeleteMii)
		auth.GET("/artisans", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/panel/artisans/1")
		})
		auth.GET("/artisans/:page", panel.ViewArtisans)
		auth.POST("/artisans/search", panel.SearchArtisans)
		auth.GET("/artisans/details/:artisan_id", panel.ViewArtisanDetails)
	}

	r.POST("/cgi-bin/conpost.cgi", conPost)
	r.POST("/cgi-bin/convote2.cgi", conVote)
	r.POST("/cgi-bin/conresult.cgi", conResult)
	r.POST("/cgi-bin/mister.cgi", mister)
	r.POST("/cgi-bin/check.cgi", check)
	r.POST("/cgi-bin/post.cgi", post)
	r.POST("/cgi-bin/vote.cgi", vote)
	r.GET("/cgi-bin/info.cgi", info)
	r.GET("/cgi-bin/ownsearch.cgi", ownSearch)
	r.GET("/cgi-bin/search.cgi", search)

	// Mii Renderer stuff
	r.POST("/cgi-bin/studio.cgi", mii.Studio)
	r.GET("/cgi-bin/render.cgi", mii.Render)

	fmt.Printf("Starting HTTP connection (%s)...\nNot using the usual port for HTTP?\nBe sure to use a proxy, otherwise the Wii can't connect!\n", config.Address)
	log.Fatalln(r.Run(config.Address))
}
