// use build constraints to work around http://code.google.com/p/go/issues/detail?id=4210
// REMOVE THIS TEMPORARILY TO RUN 'godep save'
// +build heroku

// note: need at least one blank line after build constraint
package main

import _ "bitbucket.org/liamstask/goose/cmd/goose"
