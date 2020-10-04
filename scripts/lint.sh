#!/bin/bash

go vet \
    ./cmd/* \
    ./pkg/*

golint \
    ./cmd/* \
    ./pkg/*
