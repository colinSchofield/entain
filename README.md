## Entain Technical Test

### Synopsis

This document discusses my thought process during the completion of this test, any assumptions that I made, along with a brief summary of the technical challenges and how I overcame them.

I decided to move the existing test description into a separate file, which can be found [here](TEST.md).

### Technical Challenge (*in ~~five~~ six parts*)
These sections are all included in the public git repo as separate pull requests, each following on from each other.
#### Task 0 - Initial Refactoring
When I looked at the code for this Racing API, I noticed that there were a few things missing. I decided that, as per the requirements, I needed to take *ownership* of the code base to bring it into 'a real-world environment'. I also thought about it from a *DevOps* perspective, considering  CI/CD and monitoring. By fixing some of these issues initially, the time taken will add quality to the application, providing useful scaffolding that can be of benefit to future Devs.

**Makefile**

As Go does not have any kind of a build automation tool, the UNIX Make provides us with a great starting point. 

Execute `make` or `make help` for a list of available tasks. As shown below:

```shell
Usage: make [TARGET]
Targets:
  help        Show this help.
  init        Download and install the protobuf/grpc support files.
  clean       Removes any transient build artifacts.
  generate    Generate the protobuf and gRPC Stubs & Skeletons.
  fmt         Format the Go source code.
  lint        Run lint checks.
  test        Test and Code Coverage.
  build       Build binaries on the local machine.
  docker      Build Docker images.
  run         Bring up the Racing API using Docker Compose.
```

I would say that most users might need to run `make init` to download and setup the support files (such as `protoc`). If you are a seasoned Gopher, with any kind of exposure to protobuf, this can most likely be skipped.

To run the application in a Docker container use `make run`. I considered creating a `manifest.yml` file for K8s, but there was no requirement to do this, so I left it out.

**Testing the Racing API**

Testing the application at an integration level, will require you to fire up [Postman](https://www.postman.com/) and then use the provided postman collection.

As I complete each of the Tasks, I intend to update this collection, providing a multitude of different test cases (i.e. Happy Path, alternates, plus maybe a few pen tests thrown in, for [SQL Injection](https://owasp.org/www-community/attacks/SQL_Injection) 🦹‍♂️) although I think we are mostly Okay there.

I considered including a set of tests which might be run via **Newman**, along the lines of an [API Contract test](https://medium.com/velotio-perspectives/api-testing-using-postman-and-newman-6c68c33303fc), but I could be getting a little carried away here 🙉.

**Unit Testing and TDD**

The other major thing missing is any kind of a **Unit Test** and I know for a fact, that this is *super* important. So I chose to remedy this, with a test of the package `db` via the file `races.go`. My intention is to bring this source files code coverage up to 80%. In this case I will be looking mainly at the provided filter: `MeetingIds`. I will use [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) to implement this. 

With this test code in place, it becomes easier to add new scenarios as I work on subsequent challenges.

**GitLab CI/CD**

At Entain I see that you use GitLab's CI/CD, included by the developer (i.e. `.gitlab-ci.yml`). I am a bit of a fan of [GitHub Actions](https://docs.github.com/en/actions), but I will leave out the creation of any kind of build workflow, for this challenge.

**Logging**

The Racing API does not contain any real Logging (either debug or info). This can be a great help when debugging issues with the application, when it is live in production. I intend to use the popular log provider [Logrus](github.com/sirupsen/logrus). 
This can later be hooked up into AWS CloudWatch, or ingested into another monitoring and observability tool.

#### Task 1 - Add Visible Filter

This task required that I add a boolean variable into the `ListRacesRequestFilter`, do a `make generate` and adjust the applyFilter method to modify the SQL query. If the filter is set to `false` or *no attribute* is given, ListRaces will display all races, regardless of their visibility.

There are a few new test cases to cover this scenario, plus I have added in some integration tests into the Postman collection. I also added the `racing` service into the coverage. For this sort of thing, I might use [mockery](https://medium.com/yemeksepeti-teknoloji/mocking-an-interface-using-mockery-in-go-afbcb83cc773) to generate the mocks, but this is a very straightforward interface, so I will just do it by hand.

A added in some logging  where I thought it was appropriate. Actually, it occurs to me that we probably *should not* include the proto related files (i.e. grpc.pb, pb) in the git repo, as these should always be generated each time. I am removing them from this commit.

#### Task 2 - Sort and Order by

I completed the task, including the *'Bonus points'* for ORDER/SORT-BY. The Postman collection has again been updated and I did a multi test which includes the range of new scenarios.

This one contained a potential SQL Injection threat, which I avoided by creating a map and only allowing the caller to specify from the given set of attribute names. The default, should they choose outside of this range, is the `advertised_start_time`. 

I decided to mask the names of our tables though and use the JSON names, which seems more user friendly. 