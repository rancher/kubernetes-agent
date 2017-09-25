package hostlabels

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	cache "github.com/patrickmn/go-cache"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

var (
	metadataHandler *fakeMetadataHandler
	kubeHandler     *fakeKubeNodeHandler
)

const (
	fakeMetadataURL = "http://0.0.0.0:42500/2015-12-19"
	kubeURL         = "http://0.0.0.0:42501"
)

type fakeMetadataHandler struct {
	hosts []metadata.Host
}

func (f *fakeMetadataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hb, _ := json.Marshal(f.hosts)
	w.Write(hb)
}

type fakeKubeNodeHandler struct {
	nodes map[string]*v1.Node
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (f *fakeKubeNodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// GET Node
	if r.Method == http.MethodGet {
		pathArray := strings.Split(r.URL.Path, "/")
		name := pathArray[len(pathArray)-1]
		w.Header().Set("Content-Type", "application/json")
		hb, _ := json.Marshal(f.nodes[name])
		w.Write(hb)
	}
	// Replace Node
	if r.Method == http.MethodPut {
		node := &v1.Node{}
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, node)
		f.nodes[node.Name] = node
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func TestMain(m *testing.M) {
	metadataHandler = &fakeMetadataHandler{
		hosts: []metadata.Host{},
	}

	kubeHandler = &fakeKubeNodeHandler{
		nodes: map[string]*v1.Node{},
	}

	metadataMux := http.NewServeMux()
	srv := http.Server{
		Addr:    "0.0.0.0:42500",
		Handler: metadataMux,
	}
	errChan := make(chan error, 1)
	metadataMux.Handle("/2015-12-19/hosts/", metadataHandler)
	srvLn, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalf("Error listening on tcp port 42500: %v", err)
	}
	go func() {
		errChan <- srv.Serve(tcpKeepAliveListener{srvLn.(*net.TCPListener)})
	}()
	kubeMux := http.NewServeMux()
	ksrv := http.Server{
		Addr:    "0.0.0.0:42501",
		Handler: kubeMux,
	}
	kubeMux.Handle("/api/v1/nodes/", kubeHandler)
	ksrvLn, err := net.Listen("tcp", ksrv.Addr)
	if err != nil {
		log.Fatalf("Error listening on tcp port 42501: %v", err)
	}
	go func() {
		errChan <- ksrv.Serve(tcpKeepAliveListener{ksrvLn.(*net.TCPListener)})
	}()
	intChan := make(chan int, 1)
	go func() {
		intChan <- m.Run()
	}()
	var returnVal int
	select {
	case err := <-errChan:
		log.Fatalf("Error while running metadata/kuberserver, [%v]", err)
	case returnVal = <-intChan:
	}
	os.Exit(returnVal)
}

func TestDetectsRemoval(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := kubernetesclient.NewClient(kubeURL, false)
	c := cache.New(1*time.Minute, 1*time.Minute)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test1",
			Hostname: "test1",
			Labels:   map[string]string{},
		},
	}

	kubeHandler.nodes["test1"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"test1": "val1",
			},
			Annotations: map[string]string{
				"io.rancher.labels.test1": "",
			},
			Name: "test1",
		},
	}

	sync(kubeClient, metadataClient, c)

	if _, ok := kubeHandler.nodes["test1"].Labels["test1"]; ok {
		t.Error("Label test1 was not detected as removed")
	}
}

func TestDetectsAddition(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := kubernetesclient.NewClient(kubeURL, false)
	c := cache.New(1*time.Minute, 1*time.Minute)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test2",
			Hostname: "test2",
			Labels: map[string]string{
				"test2": "val2",
			},
		},
	}

	kubeHandler.nodes["test2"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"io.kubernetes.meta": "kube.val",
			},
			Annotations: map[string]string{
				"io.kube.test": "val",
			},
			Name: "test2",
		},
	}

	sync(kubeClient, metadataClient, c)

	if _, ok := kubeHandler.nodes["test2"].Labels["test2"]; !ok {
		t.Error("Label test2 was not detected as added")
	}

	if _, ok := kubeHandler.nodes["test2"].Annotations["io.rancher.labels.test2"]; !ok {
		t.Error("Annotation was not set on addition of new label")
	}
}

func TestDetectsChange(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := kubernetesclient.NewClient(kubeURL, false)
	c := cache.New(1*time.Minute, 1*time.Minute)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test3",
			Hostname: "test3",
			Labels: map[string]string{
				"test3": "val3",
			},
		},
	}

	kubeHandler.nodes["test3"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"test3": "valx",
			},
			Annotations: map[string]string{
				"io.rancher.labels.test3": "",
			},
			Name: "test3",
		},
	}

	sync(kubeClient, metadataClient, c)

	if val := kubeHandler.nodes["test3"].Labels["test3"]; val != "val3" {
		t.Error("Label test3 was not detected as changed")
	}

	if _, ok := kubeHandler.nodes["test3"].Annotations["io.rancher.labels.test3"]; !ok {
		t.Error("Annotation was not set on addition of new label")
	}
}
