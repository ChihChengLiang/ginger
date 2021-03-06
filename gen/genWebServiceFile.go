/*
 * Copyright 2015 Ming-der Wang<ming@log4analytics.com> All right reserved.
 * Licensed by MIT License
 */
package gen

import (
	//"fmt"
	"os"
	"text/template"
)

var (
	gingerTemplate = template.Must(template.New("ginger").Parse(
		`// generated by ginger from go generate -- DO NOT EDIT

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
        "github.com/tommy351/gin-cors"
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
	db, err := gorm.Open("sqlite3", cfg.DbName)
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
	s.Migrate(cfg)
        db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	{{.VariableName}}Resource := &{{.TypeName}}Resource{db: db}

	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
        r.Use(cors.Middleware(cors.Options{}))

	r.GET("/{{.VariableName}}", {{.VariableName}}Resource.GetAll{{.TypeName}}s)
	r.GET("/{{.VariableName}}/:id", {{.VariableName}}Resource.Get{{.TypeName}})
	r.POST("/{{.VariableName}}", {{.VariableName}}Resource.Create{{.TypeName}})
	r.PUT("/{{.VariableName}}/:id", {{.VariableName}}Resource.Update{{.TypeName}})
	r.PATCH("/{{.VariableName}}/:id", {{.VariableName}}Resource.Patch{{.TypeName}})
	r.DELETE("/{{.VariableName}}/:id", {{.VariableName}}Resource.Delete{{.TypeName}})

	r.Run(cfg.SvcHost)

	return nil
}
{{end}}
`))
)

func GenWebService(path string) {
	types := findTypes(path)
	output, err := os.OpenFile("web_service.go", os.O_WRONLY|os.O_CREATE, 0600)
	defer output.Close()
	checkError(err, "could not open output file")
	gingerTemplate.Execute(output, AllType{types})
}
