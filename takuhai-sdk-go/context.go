package takuhai_sdk_go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
	"log"
	"net/http"
	"net/url"
)

type Context struct {
	Context        context.Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	managerAddr    *url.URL
	workflowID     string
	stepID         string
	jobID          string
}
type JobInfo struct {
	Time time.Duration
	RAM float64
	CPU float64
	
}//	 	ジョブ側からノード側に渡す実行時間などの情報
type Content struct {
	Body           []byte        
	Runtime        time.Duration
	RAM         float64
	CPU         float64
}

func (c *Context) Bind(v interface{}) error {
	if err := json.NewDecoder(c.Request.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (c *Context) Next(body []byte, m JobInfo) {
	cli := http.DefaultClient
	u := *c.managerAddr
	u.Path = fmt.Sprintf("/workflows/%s/steps/%s/next", c.workflowID, c.stepID)
	m.RAM = m.RAM
	log.Printf("context runtime") 
	log.Println(m.Time)
	log.Printf("context ram")
	log.Println(m.RAM)
	log.Printf("context cpu")
	log.Println(m.CPU)
	info := &Content {
	Body: body,
	Runtime: m.Time,
	RAM: m.RAM,
	CPU: m.CPU,
	}
		reqBody, err := json.Marshal(&info)
		
		if err != nil {
		log.Println(err)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(reqBody))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("takuhai-job-id", c.jobID)
	if _, err := cli.Do(req); err != nil {
		log.Println(err)
		return
	}
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Write(body)
}

func (c *Context) Finish() {
	cli := http.DefaultClient
	u := *c.managerAddr
	u.Path = fmt.Sprintf("/workflows/%s/steps/%s/finish", c.workflowID, c.stepID)
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("takuhai-job-id", c.jobID)
	if _, err := cli.Do(req); err != nil {
		log.Println(err)
		return
	}
	c.ResponseWriter.WriteHeader(http.StatusOK)
}

func (c *Context) Fail(body []byte) {
	cli := http.DefaultClient
	u := *c.managerAddr
	u.Path = fmt.Sprintf("/workflows/%s/steps/%s/fail", c.workflowID, c.stepID)
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("takuhai-job-id", c.jobID)
	if _, err := cli.Do(req); err != nil {
		log.Println(err)
		return
	}
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Write(body)
	c.Context.Done()
}
