# srechallenge
This is a golang program that takes a desired config as a yaml file and a current configuration as a json file, and returns a list of actions need for bringing the current config in-line with the desired configuration.

## Requirements
- The [Go](https://go.dev/doc/install) programming language.  
- This repo locally. Run `git clone https://github.com/ashvarts/sre-challenge.git`


## Usage
```
Usage of ./srechallenge:
  -current-config string
        the path to the current config (json api result)
  -desired-config string
        the path to the desired config (yaml configuration)
```

## How to Run
A `config.yaml` and `input.json` files are provided and represent the desired and current configurations.
1. To run the application, clone this repository.
2. Run `go run srechallenge -current-config input.json -desired-config config.yaml`.  

OR
1. Compile a binary for your system by running `go build`, this will create a binary called `srechallenge` in your local directory.
2. Run `./srechallenge -current-config input.json -desired-config config.yaml`.   

## Testing
Tests are included in main_test.go. To run the test clone this repo, and run: `go test -v`.

## Question
> Providing that the external API may be unreliable and the alert configurations can be changed by sneaky SREs manually, how would you design such an application to ensure there’s no accumulated drift between the real alerts and the desired alerts configuration?

A simple solution would involve running this applicaiton periodically, either in a loop with a sensible sleep setting or with a cron job. However, this could result in a delay between a change in configuration and application. A more sophisticated solution could also add a trigger/webhook that would fire when a new configuration is submitted and would trigger the application for immediate reconciliation. 

## Todo 
- [ ] Refactor tests
    - add more test cases to cover missing scenarios.
    - refactor code into smaller functions.
    - test "Summary".
- [ ] Refactor "Reconcile" and "Main" into smaller pieces. 
    - Add logic to delete everything and return if desired config is empty.
- [ ] Add/improve validation for inputs.