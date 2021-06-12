# Contributing

## Styling

1. For style, there is no personal preference, always use `go fmt` as a standard.


## Principles:


1. **DRY:** Always use variables, if some piece of info is occuring more than once. Define a function if a piece of code is repeating more than once.

1. **YAGNI:** No need to plan for future, if we ever need it we will implement it later..

1. **Do not hardcode:** use passable methods (flags, config, env) for configurations that can be changed in different deployments/setups and do not hardcode values in code.

1. **Write tests:** Please do not merge any feature into master unless you have writen tests for it.


## Logging

[klog](https://github.com/kubernetes/klog) is used in this repo for logging, and not any other logging tool should be used (logrus, etc). Also `fmt` stdout/stderr should be avoided.

Always use V levels for your log, unless:

* It is a fatal error which causes process to be terminated

* It is only printed once on the startup/termination, such as config load info, driver info, etc.

Log levels:

1. Errors: An operation did not completed and the component didn't have the expected functionality. However this error does not termiante the app process.

2. Warnings: Operation completed, but there was some issues.

3. Info: Important events:

    * gRPC Service.Method calls

4. Additional info for events:

    * Operation info logs in grpc methods such as resource creation, deletion, update, etc.

5. Debugging:

    * gRPC request and response logs

6. Debugging:

    * HTTP URL/METHOD requests sent to external components

7. Debugging:

    * HTTP Headers, Body requests sent to external components
