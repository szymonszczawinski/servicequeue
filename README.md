# servicequeue



## About

Simple implementation of a general usage service than can execute its methods on a separated goroutine using queue of executable jobs.

Each service has to embedd Service base struct

Each service (e.g. dummyService) such service methods has to be exposed via its own interface
that embedds IProxy interafce (to have common one for registration purpose).

Each service has to have its own proxy defined that implements service interface and embedds JobProxy struct and has interface that points to actual implementation.

Each service has to have Start method implemented in which method Start from base service is called and service is also registered ( as its own proxy )
in ServiceProvider.

When service is got from ServiceProvider it is in fact a proxy.

Inside proxy implementation of a method, a Job has to be created that executes target implementation of a method.
For void methods its rather simple, for non-void methods result channel has to be introduced to get method execution result.


