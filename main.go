/*
 * MIT License
 * Copyright 2015 Ming-der Wang<ming@log4analytics.com>
 */
package main

import (
	"fmt"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
)

type GenType struct {
	TypeName     string // could be "Slack"
	VariableName string // could be "slack"
}

type AllType struct {
	Types []GenType
}

var (
	gingerTemplate = template.Must(template.New("ginger").Parse(`
// generated by ginger from go generate -- DO NOT EDIT

package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
	Token      string
	Url        string
}

{{range .Types}}
type {{.TypeName}}Service struct {
}

func (s *{{.TypeName}}Service) getDb(cfg Config) (gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "/tmp/"+cfg.DbName)
	//db.LogMode(true)
	return db, err
}

func (s *{{.TypeName}}Service) Migrate(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	db.AutoMigrate(&{{.TypeName}}{})
	return nil
}
func (s *{{.TypeName}}Service) Run(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	slackResource := &{{.TypeName}}Resource{db: db}

	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	r.GET("/{{.VariableName}}", slackResource.GetAll{{.TypeName}}s)
	r.GET("/{{.VariableName}}/:id", slackResource.Get{{.TypeName}})
	r.POST("/{{.VariableName}}", slackResource.Create{{.TypeName}})
	r.PUT("/{{.VariableName}}/:id", slackResource.Update{{.TypeName}})
	r.PATCH("/{{.VariableName}}/:id", slackResource.Patch{{.TypeName}})
	r.DELETE("/{{.VariableName}}/:id", slackResource.Delete{{.TypeName}})

	r.Run(cfg.SvcHost)

	return nil
}
{{end}}
)
`))
)

func init() {
	log.SetFormatter(&log.TextFormatter{}) // or JsonFormatter
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
}

func checkError(err error) {
	if err == nil {
		log.WithFields(log.Fields{
			"file": "main.go",
		}).Fatal("test")
		panic(err)
	}
}

func main() {
	checkError(nil)
	fmt.Println()
}
