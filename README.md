<br><br><br><br><br><p align="center">
  <img alt="PROBE" src="https://github.com/linyows/probe/blob/main/misc/probe.svg" width="200">
</p><br><br><br><br><br>

<p align="center">
  <a href="https://github.com/linyows/probe/actions/workflows/build.yml">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/linyows/probe/build.yml?branch=main&style=for-the-badge&labelColor=000000">
  </a>
  <a href="https://github.com/linyows/probe/releases">
    <img src="http://img.shields.io/github/release/linyows/probe.svg?style=for-the-badge&labelColor=000000" alt="GitHub Release">
  </a>
  <a href="http://godoc.org/github.com/linyows/probe">
    <img src="http://img.shields.io/badge/go-documentation-blue.svg?style=for-the-badge&labelColor=000000" alt="Go Documentation">
  </a>
</p>

Probe is a YAML-based workflow automation tool. It uses plugin-based actions to execute workflows, making it highly flexible and extensible.

Example using REST API:

```yaml
name: Example http workflow
jobs:
- name: Request REST API
  defaults:
    http:
      url: http://localhost:9000
      headers:
        authorization: Bearer {env.TOKEN}
        accept: application/json
  steps:
  - name: Get a user information
    uses: http
    with:
      get: /api/v1/me
    test: res.status == 200 && res.body.uname == foobar
  - name: Update user
    uses: http
    with:
      put: /api/v1/users/{steps[0].res.body.uid}
      body:
        profile: "I'm a software engineer living in Fukuoka."
    test: res.status == 201
```

Example of sending repeated emails:

```yaml
name: Send queue congestion experiment
jobs:
- name: Normal sender
  id: normal-sender
  repeat:
    count: 60
    interval: 10
  steps:
  - use: smtp
    with:
      addr: localhost:5871
      from: alice@msa1.local
      to: bob@mx1.local
      my-hostname: msa1-local
      subject: Experiment A
- name: Throttled sender
  id: throtteled-sender
  repeat:
    count: 60
    interval: 10
  steps:
  - use: smtp
    with:
      addr: localhost:5872
      from: carol@msa2.local
      to: bob@mx2.local
      my-hostname: msa2-local
      subject: Experiment B
- name: Export latency as CSV
  needs:
  - normal-sender
  - throtteled-sender
  waitif: sh(postqueue -p 2> /dev/null | grep -c '^[A-F0-9]') != "0"
  steps:
  - use: mail-latency
    with:
      spath: /home/vmail
      dpath: ./mail-latency.csv
```

Features
--

A probe workflow consists of jobs and steps contained in the jobs. Multiple jobs are executed asynchronously, and steps are executed in sequence. Step execution results are logged, and can be expanded in YAML using curly braces.

- Workflows can be automated using built-in http, mail, and shell actions
- Custom actions that meet your use cases can be created using protocol buffers
- Protocol-based YAML definitions provide low learning costs and high visibility

Install
--

Installation via various package managers is not yet supported, but will be soon.

```sh
go install github.com/linyows/probe/cmd/probe@latest
```

Usage
--

Run the workflow by passing the path to the yaml file where the workflow is defined to the workflow option.

```sh
probe --workflow ./worflow.yml
```

To-Do
--

Here are some additional features I'm considering:

- [ ] Support waitif and needs params in job
- [ ] Support rich output
- [ ] Support multipart/form-data in http actions
- [ ] Support some actions:
    - [ ] grpc actions
    - [ ] graphql actions
    - [ ] ssh actions
    - [ ] amqp actions
    - [ ] imap actions
    - [ ] udp actions
- [ ] Support post-actions
- [ ] Support pre-job and post-job

Author
--

[linyows](https://github.com/linyows)
