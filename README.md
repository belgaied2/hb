# HB a command line tool for HobbyFarm Administrators
HobbyFarm is a Learning Management Platform used by Rancher Labs to provide Hands-On Training and Workshop. More on HobbyFarm [here](https://github.com/hobbyfarm).
HB is a CLI tool designed to help HobbyFarm Administrators interact with HobbyFarm.
The first feature it has is downloading.

Usage:
```
hb dl --email admin@rancher.com --password mypassword --region eu1 --scenario "Name of scenario"
```
where:

* `email` and `password` flags are the admin HobbyFarm credentials
* `region`: the HobbyFarm's region in the URL: https://admin.na2.hobbyfarm.io/
* `scenario`: is the name of the scenario to download

To get inline help: 
```
hb help
```
WARNING: needs improvement to document help for subcommands.

In order to compile the code, you need to have golang installed on your system and, after cloning, run :
```
go build hb.go scenario.go
```


