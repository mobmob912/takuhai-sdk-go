package takuhai_sdk_go

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Client interface {
	Run(handler func(c *Context)) error
}

type client struct {
	id          string
	workflowID  string
	stepID      string
	port        string
	handler     func(c *Context)
	managerAddr *url.URL
}

func NewClient() Client {
	return &client{}
}

func (e *client) getEnvs() error {
	managerAddr := os.Getenv("managerAddr")
	if managerAddr == "" {
		return errors.New("manager addr is missing")
	}
	ma, err := url.Parse(fmt.Sprintf("http://%s:2317", managerAddr))
	if err != nil {
		return err
	}
	e.managerAddr = ma
	port := os.Getenv("takuhaiJobPort")
	if port == "" {
		return errors.New("takuhaiJobPort is missing")
	}
	e.port = port
	workflowID := os.Getenv("workflowID")
	if workflowID == "" {
		return errors.New("workflowID is missing")
	}
	e.workflowID = workflowID
	stepID := os.Getenv("stepID")
	if stepID == "" {
		return errors.New("stepID id is missing")
	}
	e.stepID = stepID
	return nil
}

func (e *client) Run(handler func(c *Context)) error {
	if err := e.getEnvs(); err != nil {
		return err
	}
	e.handler = handler
	http.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
		parentCtx, cancel := context.WithCancel(r.Context())
		jobID := r.Header.Get("takuhai-job-id")
		ctx := &Context{
			Context:        parentCtx,
			Request:        r,
			ResponseWriter: w,
			managerAddr:    e.managerAddr,
			workflowID:     e.workflowID,
			stepID:         e.stepID,
			jobID:          jobID,
		}
		e.handler(ctx)
		cancel()
	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	if err := e.NotifyReady(); err != nil {
		return err
	}
	log.Println("serving... port:" + e.port)
	return http.ListenAndServe(":"+e.port, nil)
}

func (e *client) NotifyReady() error {
	c := http.DefaultClient
	u := *e.managerAddr
	u.Path = fmt.Sprintf("/workflows/%s/steps/%s", e.workflowID, e.stepID)
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}
	if _, err := c.Do(req); err != nil {
		return err
	}
	return nil
}

//
//func (e *client) getJobID() error {
//	req, err := http.NewRequest("POST", "/jobs/"+e.jobName, nil)
//	if err != nil {
//		return err
//	}
//	if err := req.Write(e.managerConn); err != nil {
//		return err
//	}
//	resp, err := http.ReadResponse(bufio.NewReader(e.managerConn), req)
//	if err != nil {
//		return err
//	}
//	var respStruct struct {
//		ID string `json:"id"`
//	}
//	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
//		return err
//	}
//	e.id = respStruct.ID
//	log.Println(e.id)
//	return nil
//}
//
