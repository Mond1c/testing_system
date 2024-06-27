// Package internal contains internal logic of the application.
package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// TODO: Think about this

// Task is an interface that can be executed by local Worker or remote Worker.
type Task interface {
	Local() error
	Remote(addr string) error
}

// RunTask is an implementation of Task interface that needed to run tests.
type RunTask struct {
	duration           int64
	startTime          string
	testPathForProblem string
	testCount          int
	username           string
	compileFile        string
	run                *Run
	action             func(*RunInfo)
}

// CreateRunTask creates RunTask for running tests
func CreateRunTask(duration int64, startTime, testPathForProblem string, testCount int, username, compileFile string,
	run *Run, action func(*RunInfo)) *RunTask {
	return &RunTask{
		duration:           duration,
		startTime:          startTime,
		testPathForProblem: testPathForProblem,
		testCount:          testCount,
		username:           username,
		compileFile:        compileFile,
		run:                run,
		action:             action,
	}
}

// Local runs tests locally
func (rt *RunTask) Local() error {
	ts := NewTestingSystem(rt.duration, rt.startTime, rt.testPathForProblem, rt.testCount, rt.username, rt.compileFile)
	info, err := ts.RunTests(rt.run)
	if err != nil {
		return err
	}
	rt.action(info)
	return nil
}

// mustOpen opens file or returns error
func mustOpen(path string) (*os.File, error) {
	r, err := os.Open(path)
	return r, err
}

// Remote sends post request to remote Worker with the specified address and wait for response.
func (rt *RunTask) Remote(addr string) error {
	httpClient := &http.Client{}
	file, err := mustOpen(rt.run.fileName)
	if file == nil {
		return err
	}
	values := map[string]io.Reader{
		"file":      file,
		"language":  strings.NewReader(rt.run.language),
		"problem":   strings.NewReader(rt.run.problem),
		"username":  strings.NewReader(rt.run.username),
		"duration":  strings.NewReader(strconv.FormatInt(rt.duration, 10)),
		"startTime": strings.NewReader(rt.startTime),
		"userId":    strings.NewReader(rt.run.userId),
		"filename":  strings.NewReader(rt.run.fileName),
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return err
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}
	}
	w.Close()
	req, err := http.NewRequest("POST", addr, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}
	var runInfo *RunInfo
	err = json.NewDecoder(res.Body).Decode(&runInfo)
	if err != nil {
		return err
	}
	rt.action(runInfo)
	return nil
}

// CreateRejudgeTask creates Task for rejudging tests
// FIXME: Implement me
func CreateRejudgeTask(runInfo *RunInfo) Task {
	panic("implement me")
}

// TestingQueue is a queue for running tasks
type TestingQueue struct {
	queue   chan Task
	workers chan Worker
}

// NewTestingQueue creates new pointer to TestingQueue
func NewTestingQueue(workers chan Worker) *TestingQueue {
	return &TestingQueue{
		// FIXME: Think about buffer size
		queue:   make(chan Task, 1000),
		workers: workers,
	}
}

// PushTask pushes task to the queue
func (tq *TestingQueue) PushTask(task Task) {
	tq.queue <- task
}

// Update updates the queue
func (tq *TestingQueue) Update() {
	for {
		task := <-tq.queue
		go func() {
			// Deadlock?
			worker := <-tq.workers
			defer func() { tq.workers <- worker }()
			err := worker.RunTask(task)
			if err != nil {
				log.Print(err)
			}
		}()
	}
}

// CreateSimpleTestingQueue creates simple TestingQueue
func CreateSimpleTestingQueue() *TestingQueue {
	workers := make(chan Worker, 10)
	for i := 0; i < 10; i++ {
		worker := NewRemoteWorker("http://127.0.0.1:8081/test")
		workers <- worker
	}
	return NewTestingQueue(workers)
}

var MyTestingQueue = CreateSimpleTestingQueue()
