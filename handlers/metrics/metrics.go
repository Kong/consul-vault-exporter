package metrics

import (
	"bytes"
	//	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

	if os.Getenv("CONSUL_ADDRESS") == "" {
		statusUrl = "consul:8500"
	} else {
		statusUrl = os.Getenv("CONSUL_ADDRESS")
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

	nodes, err := DiscoverNodes(c, url)
	if err != nil {
		return metricString, err
	}
	for _, n := range nodes {
		statusClient := http.Client{
						Timeout: time.Second * 3,
		}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/", n), nil)
		if err != nil {
						return metricString, err
		}
		res, getErr := statusClient.Do(req)
		if getErr != nil {
						handleError(c, getErr)
						return metricString, getErr
		}
		if res.StatusCode != 200 {
						c.JSON(500, gin.H{"message": "Call to the api endpoint failed", "http_status": res.StatusCode})
						return metricString, errors.New("HTTP Response: was not 200")
		}
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
						return metricString, readErr
		}
		fmt.Println(body)
	}

	metricString.WriteString(fmt.Sprintf("# HELP Can we fetch the catalog from consul\n# TYPE consul_catalog_available gauge\nconsul_catalog_available %d\n", 1))
	return metricString, nil
}

func DiscoverNodes(c *gin.Context, url string) ([]string, error) {
	var nodeList []string

	consulConf := api.Config{Address: url}
	client, err := api.NewClient(&consulConf)
	if err != nil {
		return nodeList, err
	}

	catalog := client.Catalog()
	srvs := api.QueryOptions{}
	//res, _, _ := catalog.Services(&srvs)
	res, _, err := catalog.Service("vault", "", &srvs)
	if err != nil {
		return nodeList, err
	}

	r, _ := regexp.Compile(`^vault-.*$`)
	for _, s := range res {
			if r.MatchString(s.Node) == true {
				nodeList = append(nodeList, fmt.Sprintf("%s:%d", s.Node, s.ServicePort))
			}
	}
	return nodeList, nil
}
