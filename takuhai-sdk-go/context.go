package takuhai_sdk_go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (c *Context) Bind(v interface{}) error {
	if err := json.NewDecoder(c.Request.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (c *Context) Next(body []byte) {
	cli := http.DefaultClient
	u := *c.managerAddr
	u.Path = fmt.Sprintf("/workflows/%s/steps/%s/next", c.workflowID, c.stepID)
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
