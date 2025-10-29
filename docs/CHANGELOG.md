# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.0]
### [Added]
- Implemented JSON based Get and Put properties for WebDAV (PROPFIND and PROPPATCH methods)

## [0.3.0]
### [Internal]
- CI/CD with Github Actions to build and deploy Docker image to Docker Hub on new release tag
- Updated Dockerfile

### [Added]
- TODO: Add changelog entries here for the next release
- Forgot to do that when fighting with Github actions deploys

## [0.2.0]
### [Added]
- Digest Authentication middleware implementation (RFC 2069)
- WebDAV compatibility for Digest Authentication
- Implemented `/users/me` returning currently logged-in user
- Customized logging middleware and created custom logger
- Implemented basics for webdav support
- Implemented basic file operations: list files, create directory, upload file, download file, delete file

## [0.1.0]
### [Internal]
- Local deployment docker-compose

### [Added]
 - Basic API project structure
 - Healthcheck endpoint
 - Variables from .env support
 - Basic DB connection file structure
 - JSON to struct data parser
 - Auth endpoints to create and authorize user
 - Custom HttpError to standardize api error response in JSON