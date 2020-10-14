#!/bin/bash

go vet . ./cmd

golint . ./cmd
