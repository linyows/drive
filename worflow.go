package probe

import (
	"fmt"
	"sync"
	"time"
)

type Workflow struct {
	Name string `yaml:"name",validate:"required"`
	Jobs []Job  `yaml:"jobs",validate:"required"`
}

func (w *Workflow) Start() {
	ctx := w.createContext()
	var wg sync.WaitGroup

	for _, job := range w.Jobs {
		// No repeat
		if job.Repeat == nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				job.Start(ctx)
			}()
			continue
		}

		// Repeat
		for i := 0; i < job.Repeat.Count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				job.Start(ctx)
			}()
			time.Sleep(time.Duration(job.Repeat.Interval) * time.Second)
		}
	}

	wg.Wait()
}

func (w *Workflow) createContext() JobContext {
	return JobContext{
		Envs: getEnvMap(),
		Logs: []map[string]any{},
	}
}

type JobContext struct {
	Envs map[string]string `expr:"env"`
	Logs []map[string]any  `expr:"steps"`
}

type Repeat struct {
	Count    int `yaml:"count",validate:"required,gte=0,lt=100"`
	Interval int `yaml:"interval,validate:"gte=0,lt=600"`
}

type Step struct {
	Name   string         `yaml:"name"`
	Uses   string         `yaml:"uses" validate:"required"`
	With   map[string]any `yaml:"with"`
	errors error
}

type Job struct {
	Name     string  `yaml:"name",validate:"required"`
	Steps    []Step  `yaml:"steps",validate:"required"`
	Repeat   *Repeat `yaml:"repeat"`
	Defaults any     `yaml:"defaults"`
	ctx      *JobContext
}

func (j *Job) Start(ctx JobContext) {
	j.ctx = &ctx

	for i, st := range j.Steps {
		expW := EvaluateExprs(st.With, ctx)
		ret, err := RunActions(st.Uses, []string{}, expW)
		if err != nil {
			st.errors = err
			continue
		}
		req, okreq := ret["req"].(map[string]any)
		res, okres := ret["res"].(map[string]any)
		if okres {
			// parse json and sets
			body, okbody := res["body"].(string)
			if okbody && isJSON(body) {
				res["bodyjson"] = mustMarshalJSON(body)
			}
		}
		if okreq && okres {
			ShowVerbose(i, req, res)
		} else {
			fmt.Printf("---------- Step %d ----------\n%#v\n", i, ret)
		}
		ctx.Logs = append(ctx.Logs, ret)
	}
}

func ShowVerbose(i int, req, res map[string]any) {
	fmt.Printf("---------- Step %d ----------\nRequest:\n", i)
	for k, v := range req {
		nested, ok := v.(map[string]any)
		if ok {
			fmt.Printf("  %s:\n", k)
			for kk, vv := range nested {
				fmt.Printf("    %s: %#v\n", kk, vv)
			}
		} else {
			fmt.Printf("  %s: %#v\n", k, v)
		}
	}
	fmt.Printf("Response:\n")
	for k, v := range res {
		nested, ok := v.(map[string]any)
		if ok {
			fmt.Printf("  %s:\n", k)
			for kk, vv := range nested {
				fmt.Printf("    %s: %#v\n", kk, vv)
			}
		} else {
			fmt.Printf("  %s: %#v\n", k, v)
		}
	}
}