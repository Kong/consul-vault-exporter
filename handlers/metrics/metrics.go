package metrics

import (
	"bytes"
	"encoding/json"
	//"errors"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	//"time"
	"github.com/hashicorp/consul/api"
	"regexp"
)

type VauthHealth struct {
	Initialized                bool   `json:"initialized"`
	Sealed                     bool   `json:"sealed"`
	Standby                    bool   `json:"standby"`
	ReplicationPerformanceMode string `json:"replication_performance_mode"`
	ReplicationDrMode          string `json:"replication_dr_mode"`
	ServerTimeUtc              int    `json:"server_time_utc"`
	Version                    string `json:"version"`
	ClusterName                string `json:"cluster_name"`
	ClusterID                  string `json:"cluster_id"`
}

func Redirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/metrics")
}

func handleError(c *gin.Context, errorString error) {
	c.JSON(500, gin.H{"message": errorString, "status": "error"})
}

func bool2float(b bool) float64 {
	if b {
		return 1
	}
	return 0
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
	var nodeCount uint64

	nodes, err := DiscoverNodes(c, url)
	if err != nil {
		return metricString, err
	}
	for _, n := range nodes {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		res, err := client.Get(fmt.Sprintf("https://%s/v1/sys/health", n))
		if err != nil {
			return metricString, err
		}
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			return metricString, readErr
		}
		vaultHealth := VauthHealth{}
		jsonErr := json.Unmarshal(body, &vaultHealth)
		if jsonErr != nil {
			return metricString, jsonErr
		}
		metricString.WriteString(fmt.Sprintf("# HELP is vault initialized?\n# TYPE vault_initialized gauge\nvault_initialized{instance=\"%s\",cluster=\"%s\",version=\"%s\"} %.f\n", n, vaultHealth.ClusterName, vaultHealth.Version, bool2float(vaultHealth.Initialized)))
		metricString.WriteString(fmt.Sprintf("# HELP is vault sealed?\n# TYPE vault_sealed gauge\nvault_sealed{instance=\"%s\",cluster=\"%s\",version=\"%s\"} %.f\n", n, vaultHealth.ClusterName, vaultHealth.Version, bool2float(vaultHealth.Sealed)))
		metricString.WriteString(fmt.Sprintf("# HELP is vault standby?\n# TYPE vault_standby gauge\nvault_standby{instance=\"%s\",cluster=\"%s\",version=\"%s\"} %.f\n", n, vaultHealth.ClusterName, vaultHealth.Version, bool2float(vaultHealth.Standby)))
		nodeCount += 1
	}

	metricString.WriteString(fmt.Sprintf("# HELP discovered node count\n# TYPE vault_node_count gauge\nvault_node_count %d\n", nodeCount))

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
	res, _, err := catalog.Service("vault", "", &srvs)
	if err != nil {
		return nodeList, err
	}

	r, _ := regexp.Compile(`^vault-.*$`)
	for _, s := range res {
		if r.MatchString(s.Node) == true {
			nodeList = append(nodeList, fmt.Sprintf("%s:%d", s.ServiceAddress, s.ServicePort))
		}
	}
	return nodeList, nil
}
