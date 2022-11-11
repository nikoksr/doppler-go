<div align="center">

&nbsp;
<h1>doppler-go</h1>
<p><i>A Go client for the <a href="https://www.doppler.com/">Doppler </a>  API.</i></p>

&nbsp;

[![codecov](https://codecov.io/gh/nikoksr/doppler-go/branch/main/graph/badge.svg?token=9KTRRRWM5A)](https://codecov.io/gh/nikoksr/doppler-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikoksr/doppler-go)](https://goreportcard.com/report/github.com/nikoksr/doppler-go)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ff90807c42154df9b12a5f03d30a7160)](https://www.codacy.com/gh/nikoksr/doppler-go/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nikoksr/doppler-go&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/8d58f3077a2b6ee2ac57/maintainability)](https://codeclimate.com/github/nikoksr/doppler-go/maintainability)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/doppler-go)
</div>

&nbsp;
## About <a id="about"></a>

doppler-go is a Go client for the [Doppler](https://www.doppler.com/) API. It provides a simple and easy to use interface for interacting with the [Doppler API](https://docs.doppler.com/reference/api).

&nbsp;
## Stability <a id="stability"></a>

This project is currently in an early stage of development. I have not yet confirmed the functionality of all endpoints through extensive end-to-end testing. Therefore, I cannot guarantee for any stability or correctness.

I plan to support this project for the foreseeable future. However, I cannot guarantee that I will be able to fix bugs or add new features in a timely manner.

If you find any bugs or have any suggestions, please open an issue or a pull request.

&nbsp;
## Features <a id="features"></a>

* Doppler REST API v3:
  * Audit
  * Auth
  * Configs
  * Config Logs
  * Dynamic Secrets
  * Environments
  * Projects
  * Secrets
  * Service Tokens
  * Token Sharing
  * Workplaces

&nbsp;
## Install <a id="install"></a>

```sh
go get -u github.com/nikoksr/doppler-go
```

&nbsp;
## Example usage <a id="usage"></a>

```go
package main

import (
  "context"
  "fmt"
  "log"

  "github.com/nikoksr/doppler-go"
  "github.com/nikoksr/doppler-go/secret"
)

func main() {
  // Set your API key
  doppler.Key = "YOUR_API_KEY"

  // List all your secrets
  secrets, err := secret.List(context.Background(), &doppler.SecretListOptions{
    Project: "YOUR_PROJECT",
    Config:  "YOUR_CONFIG",
  })
  if err != nil {
    log.Fatalf("failed to list secrets: %v", err)
  }

  for name, value := range secrets {
    fmt.Printf("%s: %v\n", name, value)
  }
}
```

&nbsp;
## Contributing <a id="contributing"></a>

Contributions of all kinds are very welcome! Feel free to check
our [open issues](https://github.com/nikoksr/doppler-go/issues). Please also take a look at
the [contribution guidelines](https://github.com/nikoksr/doppler-go/blob/main/CONTRIBUTING.md).

&nbsp;
## Show your support <a id="support"></a>

Please give a ⭐️ if you like this project! This helps us to get more visibility and helps other people to find this
project.
