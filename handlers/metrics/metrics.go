package metrics

import (
	"bytes"
	//	"encoding/json"
	//	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	//	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"github.com/hashicorp/consul/api"
)

func Redirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/metrics")
}

func handleError(c *gin.Context, errorString error) {
	c.JSON(500, gin.H{"message": errorString, "status": "error"})
}

func Metrics(c *gin.Context) {
	var statusUrl string

	if os.Getenv("CONSUL_URL") == "" {
		statusUrl = "http://consul:8500/"
	} else {
		statusUrl = os.Getenv("CONSUL_URL")
	}

	perfdata, err := ScrapeMetrics(c, statusUrl)
	if err != nil {
		handleError(c, err)
		return
	}
	outData := fmt.Sprintf("%s", perfdata.String())

	//append(perfdata, apidata)
	c.Data(200, "application/json; charset=utf-8", []byte(outData))
}

func ScrapeMetrics(c *gin.Context, url string) (bytes.Buffer, error) {
	var metricString bytes.Buffer

	consulConf := api.Config{Address: "consul-infra-aws-us-west-1.kongcloud.io:8500"}
	client, err := api.NewClient(&consulConf)
	if err != nil {
		panic(err)
	}

	catalog := client.Catalog()
	srvs := api.QueryOptions{}
	//res, _, _ := catalog.Services(&srvs)
	res, _, err := catalog.Service("vault", "", &srvs)
	if err != nil {
		panic(err)
	}

	r, _ := regexp.Compile(`^vault-.*$`)
	for _,s := range res {
			if r.MatchString(s.Node) == true {
				fmt.Println(s.Node)
				fmt.Println(s.ServicePort)
			}
	}


	metricString.WriteString(fmt.Sprintf("# HELP Can we fetch the catalog from consul\n# TYPE consul_catalog_available gauge\nconsul_catalog_available %d\n", 1))
	return metricString, nil
}
