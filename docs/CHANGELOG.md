# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### [Added]
- Digest Authentication middleware implementation (RFC 2069)
- WebDAV compatibility for Digest Authentication
- Implemented `/users/me` returning currently logged in user

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