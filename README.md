# Article Web Service

[![Echo Logo](https://cdn.labstack.com/images/echo-logo.svg)](https://echo.labstack.com)

[![Sourcegraph](https://sourcegraph.com/github.com/labstack/echo/-/badge.svg?style=flat-square)](https://sourcegraph.com/github.com/labstack/echo?badge)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/labstack/echo/v4)
[![Go Report Card](https://goreportcard.com/badge/github.com/labstack/echo?style=flat-square)](https://goreportcard.com/report/github.com/labstack/echo)
[![Build Status](http://img.shields.io/travis/labstack/echo.svg?style=flat-square)](https://travis-ci.org/labstack/echo)
[![Codecov](https://img.shields.io/codecov/c/github/labstack/echo.svg?style=flat-square)](https://codecov.io/gh/labstack/echo)
[![Forum](https://img.shields.io/badge/community-forum-00afd1.svg?style=flat-square)](https://github.com/labstack/echo/discussions)
[![Twitter](https://img.shields.io/badge/twitter-@labstack-55acee.svg?style=flat-square)](https://twitter.com/labstack)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)

## Description

This is a simple web service for managing articles. It allows you to create new articles and retrieve a list of articles. The service is built using the Go programming language, Echo framework, and PostgreSQL database.

## Pre-requisites

Before running the project, make sure you have the following dependencies installed:

- Go (1.16 or higher)
- PostgreSQL
- Docker

## Running the project

To run the project, you need to set the following environment variables:

- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `DB_USER` - PostgreSQL user
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - PostgreSQL database name
- `DOCKER_NAME` - Docker container name
- `APP_PORT` - Port on which the application will be available

After that, you can run the following command:

```bash
make run
```

## Migration

To create a database schema, you need to run the following command:

```bash
make migrate-up
```
