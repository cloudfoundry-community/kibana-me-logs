Go Cloud Foundry Environment Package (cfenvnested)
==================================================

NOTE: this is a fork of https://github.com/cloudfoundry-community/go-cfenv that allows service credentials with nested structures. One day we will figure out how to blend the two together.

### Overview

[![GoDoc](https://godoc.org/github.com/cloudfoundry-community/go-cfenvnested?status.png)](https://godoc.org/github.com/cloudfoundry-community/go-cfenvnested)

`cfenvnested` is a package to assist you in writing Go apps that run on [Cloud Foundry](http://cloudfoundry.org). It provides convenience functions and structures that map to Cloud Foundry environment variable primitives (http://docs.cloudfoundry.com/docs/using/deploying-apps/environment-variable.html).

### Build Status

-	[![Build Status - Master](https://travis-ci.org/cloudfoundry-community/go-cfenvnested.svg?branch=master)](https://travis-ci.org/cloudfoundry-community/go-cfenvnested) `Master`
-	[![Build Status - Develop](https://travis-ci.org/cloudfoundry-community/go-cfenvnested.svg?branch=develop)](https://travis-ci.org/cloudfoundry-community/go-cfenvnested) `Develop`

### Usage

`go get github.com/cloudfoundry-community/go-cfenvnested`

```go
package main

import (
	"github.com/cloudfoundry-community/go-cfenvnested"
)

func main() {
	appEnv := cfenvnested.Current()

	fmt.Println("ID:", appEnv.ID)
	fmt.Println("Index:", appEnv.Index)
	fmt.Println("Name:", appEnv.Name)
	fmt.Println("Host:", appEnv.Host)
	fmt.Println("Port:", appEnv.Port)
	fmt.Println("Version:", appEnv.Version)
	fmt.Println("Home:", appEnv.Home)
	fmt.Println("MemoryLimit:", appEnv.MemoryLimit)
	fmt.Println("WorkingDir:", appEnv.WorkingDir)
	fmt.Println("TempDir:", appEnv.TempDir)
	fmt.Println("User:", appEnv.User)
	fmt.Println("Services:", appEnv.Services)
}
```

### Contributing

Pull requests welcomed. Please ensure you make your changes in a branch off of the `develop` branch, not the `master` branch.

Remember that you probably need to patch https://github.com/cloudfoundry-community/go-cfenv project as well.
