/*
 * Copyright 2015 Ming-der Wang<ming@log4analytics.com> All right reserved.
 * Licensed by MIT License
 */
package gen

import (
	//"fmt"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
)

var (
	gingerTemplate2 = template.Must(template.New("ginger").Parse(
		`// generated by ginger from go generate -- DO NOT EDIT
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

{{range .Types}}
type {{.TypeName}}Resource struct {
	db gorm.DB
}

// @Title Create{{.TypeName}}
// @Description get string by ID
// @Accept  json
// @Param   some_id     path    int     true        "Some ID"
// @Success 201 {object} string
// @Failure 400 {object} APIError "problem decoding body"
// @Router /{{.VariableName}}/ [post]
func (tr *{{.TypeName}}Resource) Create{{.TypeName}}(c *gin.Context) {
	var {{.VariableName}} {{.TypeName}}

	if c.Bind(&{{.VariableName}}) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}
	//{{.VariableName}}.Status = {{.TypeName}}Status
	{{.VariableName}}.Created = int32(time.Now().Unix())

	tr.db.Save(&{{.VariableName}})

	c.JSON(http.StatusCreated, {{.VariableName}})
}

func (tr *{{.TypeName}}Resource) GetAll{{.TypeName}}s(c *gin.Context) {
	var {{.VariableName}}s []{{.TypeName}}

	tr.db.Order("created desc").Find(&{{.VariableName}}s)

	c.JSON(http.StatusOK, {{.VariableName}}s)
}

func (tr *{{.TypeName}}Resource) Get{{.TypeName}}(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var {{.VariableName}} {{.TypeName}}

	if tr.db.First(&{{.VariableName}}, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		c.JSON(http.StatusOK, {{.VariableName}})
	}
}

func (tr *{{.TypeName}}Resource) Update{{.TypeName}}(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var {{.VariableName}} {{.TypeName}}

	if c.Bind(&{{.VariableName}}) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}
	{{.VariableName}}.Id = int32(id)

	var existing {{.TypeName}}

	if tr.db.First(&existing, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		tr.db.Save(&{{.VariableName}})
		c.JSON(http.StatusOK, {{.VariableName}})
	}

}

func (tr *{{.TypeName}}Resource) Patch{{.TypeName}}(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	// this is a hack because Gin falsely claims my unmarshalled obj is invalid.
	// recovering from the panic and using my object that already has the json body bound to it.
	var json []Patch

	r := c.Bind(&json)
	if r != nil {
		fmt.Println(r)
	} else {
		if json[0].Op != "replace" && json[0].Path != "/status" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "PATCH support is limited and can only replace the /status path"})
			return
		}
		var {{.VariableName}} {{.TypeName}}

		if tr.db.First(&{{.VariableName}}, id).RecordNotFound() {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		} else {
			{{.VariableName}}.Status = json[0].Value

			tr.db.Save(&{{.VariableName}})
			c.JSON(http.StatusOK, {{.VariableName}})
		}
	}
}

func (tr *{{.TypeName}}Resource) Delete{{.TypeName}}(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var {{.VariableName}} {{.TypeName}}

	if tr.db.First(&{{.VariableName}}, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		tr.db.Delete(&{{.VariableName}})
		c.Data(http.StatusNoContent, "application/json", make([]byte, 0))
	}
}

func (tr *{{.TypeName}}Resource) getId(c *gin.Context) (int32, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int32(id), nil
}

/**
* on patching: http://williamdurand.fr/2014/02/14/please-do-not-patch-like-an-idiot/
 *
  * patch specification https://tools.ietf.org/html/rfc5789
   * json definition http://tools.ietf.org/html/rfc6902
*/

type Patch struct {
	Op    string ` + "`" + `json:"op" binding:"required"` + "`" + `
	From  string ` + "`" + `json:"from"` + "`" + `
	Path  string ` + "`" + `json:"path"` + "`" + `
	Value string ` + "`" + `json:"value"` + "`" + `
}
{{end}}
`))
)

func init() {
	log.SetFormatter(&log.TextFormatter{}) // or JsonFormatter
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
}

func GenResourceFile(path string) {
	types := findTypes(path)
	output, err := os.OpenFile(lowerFirstChar(path)+"_resource.go", os.O_WRONLY|os.O_CREATE, 0600)
	defer output.Close()
	checkError(err, "could not open output file")
	gingerTemplate2.Execute(output, AllType{types})
}
