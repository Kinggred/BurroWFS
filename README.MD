# BurroWFS

https://www.zotero.org/
## [Versioning](https://semver.org/)

## Development pipeline:

`main` -> `x.y-dev` -> `xyz-Short-description` -> `x.y-dev` -> `main`

#### Where:

* `main` is a project **main** branch. Code merged into it **has to** be possible to run
* `x.y-dev` (`0.1-dev`) is used to implement **major** parts of business logic. Changes merged into it **can** contain unfinished parts.
* `xyz-Short-description` (`112-PostgresDB-setup`) is used to implement **minor** business logic elements. Changes **should** be as small as possible. Number assigned to each minor branch is derived from TODO list. Commits have to be squashed on Pull.


## TODO:
It is important to keep in mind that this is not a roadmap and as such this list doesn't describe the order in which problems will be approached. 
Moreover, as this list is used to name minor changes and due to the fact that it is impossible to predict every necessary step it is destined to grow.


1. Project
   1. Initialization
      1. API
      2. Database
   2. Documentation
   3. Quality of life
      1. GitHub Actions
   4. Local Deployment

2. API:
   1. Auth
      1. Setup PostgresDB to store Users
      2. Write authorization logic
      3. [Optional] Support for Key auth ?
   2. Files
      1. Logic retrieving list of users files
      2. Logic of file upload
      3. Logic of file download
      4. Logic of file Locking
   3. WebDAV
      1. Basic WebDAV support
      2. Integrate API logic with WebDAV layer
      3. Minimal implementation to connect with file explorers
      4. Retrieving files and polishing putting files
      5. Create collections and move files around
   4. Events
      1. Handle certain S3 Events
   5. Other
      1. Error Handling
      2. Logging

3. AWS
   1. Lambda
      1. Setup necessary logic to run GO on lambda
      2. Deploy suitable Lambda with code
      3. Setup GitHub action to release
   2. Postgres
      1. Setup serverless Postgres instance
      2. Setup migrations
   3. S3
      1. Setup S3 instance
      2. Setup Events
      3. Setup Bucket integration
