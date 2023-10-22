package main

import (
    "fmt"
    "net/url"
    "gopkg.in/yaml.v3"
    "os"
    "io"
	"net/http"
    "strings"
)

type Config struct {
    Listen struct {
        Port int 
    }
    Applications []struct {
        Path string
        Forward string
        Rewrite bool `default:"false"`
    }
}



func main() {
    // read the output.yaml file 
    data, err := os.ReadFile("/config/logincrate.yml")
    
    data = []byte(os.ExpandEnv(string(data)))

    if err != nil {
        panic(err)
    }

    // create a person struct and deserialize the data into that struct
    var config Config

    err = yaml.Unmarshal([]byte(data), &config);
    if err != nil {
        fmt.Fprint(os.Stderr, "Invalid Configuration Syntax")
        fmt.Fprintf(os.Stderr, "%s", err)
        os.Exit(1)
    }

    for _, a := range config.Applications {
        _, err := url.Parse(a.Forward)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid URL ( %s ) in application with path: %s", a.Forward, a.Path)
        }
    }
    
    http.HandleFunc("/", func (rw http.ResponseWriter, req *http.Request) {
        for _, app := range config.Applications {
            if strings.HasPrefix(req.URL.Path, app.Path) {
                newUrl, _ := url.Parse(app.Forward)

                req.URL.Host = newUrl.Host
                req.URL.Scheme = newUrl.Scheme
                req.RequestURI=""

                if app.Rewrite {
                    req.URL.Path = newUrl.Path+strings.TrimPrefix(req.URL.Path, app.Path)
                }
                
                res, err := http.DefaultClient.Do(req)

                if err != nil {
                    panic(err)
                }
                for name, values := range res.Header {
                    for _, value := range values {
                        rw.Header().Set(name, value)
                    }
                }
                buf := make([]byte, 8)
                if _, err := io.CopyBuffer(rw, res.Body, buf); err != nil {
                    panic(err)
                }
                res.Body.Close()  
                return
            }
        }
        http.Error(rw, "404: Not Found", 404)
    })
    
    fmt.Printf("Listening on port %d", config.Listen.Port)
    err = http.ListenAndServe(fmt.Sprintf(":%d", config.Listen.Port), nil)
    if err != nil {
        panic(err)
    }
    
}