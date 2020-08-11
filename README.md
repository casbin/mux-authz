# Mux-authz [![Build Status](https://travis-ci.org/casbin/mux-authz.svg?branch=master)](https://travis-ci.org/casbin/mux-authz)[![Coverage Status](https://coveralls.io/repos/github/casbin/mux-authz/badge.svg?branch=master)](https://coveralls.io/github/casbin/mux-authz?branch=master)

Mux-authz is an authorization middleware for [Mux](https://github.com/gorilla/mux), it’s based on [Casbin](https://github.com/casbin/casbin). If you have better suggestions, please submit the issue.

## Installation

```
go get github.com/casbin/mux-authz
```

## Prepare

This repo is based on [Casbin](http://github.com/casbin/casbin), so you need to prepare two files in advance.

The Casbin model file describes access control models like ACL, RBAC, ABAC, etc. 

The Casbin policy file describes the authorization policy rules. 

For how to write these files, please refer to: https://github.com/casbin/casbin#get-started

## How to use

1. Create your Casbin model file [authz_model.conf](https://github.com/casbin/mux-authz/blob/master/authz_model.conf) and Casbin policy file [authz_policy.csv](https://github.com/casbin/mux-authz/blob/master/authz_policy.csv) into this folder. 

2. Load model and policy 

   ```go
   c := new(authz.CasbinAuthorizer)
   err :=c.Load("authz_model.conf", "authz_policy.csv")
   if err != nil {
   	fmt.Println(err.Error())
   }
   ```

3. Use Middleware

   ```go
   r :=mux.NewRouter()
   r.HandleFunc("/{url:[A-Za-z0-9\\/]+}", handler)
   r.Use(c.Middleware)
   ```

   Note: Now we only support check the whole path. So we recommend using path with regular expressions in the HandleFunc(). In this way, you don't have to worry about 404 due to the number of ‘/‘.For example `/book1/1` and `/bookshelf1/book1/1`.

If you have any questions, you can refer to [mux-authz_test.go](https://github.com/casbin/mux-authz/blob/master/mux-authz_test.go).

## Complete Example.

```go
package main

import (
	"fmt"
	authz "github.com/casbin/mux-authz"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	c := new(authz.CasbinAuthorizer)
	err :=c.Load("authz_model.conf", "authz_policy.csv")
	if err != nil {
		fmt.Println(err.Error())
	}

	// A very simple health check handler.
	handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	r :=mux.NewRouter()
	r.HandleFunc("/{url:[A-Za-z0-9\\/]+}", handler)
	r.Use(c.Middleware)
	log.Fatal(http.ListenAndServe(":8080",r))
}
```

Note: This plugin only supports HTTP basic authentication to get the logged-in user name, if you use other kinds of authentication like OAuth, LDAP, etc, you may need to customize this plugin.

## How to control the access

The authorization determines a request based on `{subject, object, action}`, which means what `subject` can perform what `action` on what `object`. In this plugin, the meanings are:

1. `subject`: the logged-on user name
2. `object`: the URL path for the web resource like "dataset1/item1"
3. `action`: HTTP method like GET, POST, PUT, DELETE, or the high-level actions you defined like "read-file", "write-blog"

For how to write authorization policy and other details, please refer to [the Casbin's documentation](https://github.com/casbin/casbin).

## Getting Help

- [Casbin](https://github.com/casbin/casbin)