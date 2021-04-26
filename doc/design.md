## Design decisions

1. GoLang as an implementation language

*Rationale: Current gitlab runner is implemented in GoLang, and seems to be community language of choice 
when it comes to system programming. Alternative is Rust, though with steep learning curve and lower popularity, 
it was elminated as possible choice*

2. Build runners state saved as local `json` file. Concurrency issues are avoided using poor-mans lockfile implementation. 

*Rationale: As this solution aims for MVP, any cloud based storage, like Dynamo for AWS, would require additional abstractions
and prolong implementation time. Only downside to not having persistant storage is orphaned build VMs in case of abnormal termination
of runner manager instance*

3. `ssh` is used as remote code execution protocol between gitlab runner and executor. 

4. Any provisioning parameters passed to cloud provider provisioner e.g. `RunInstances` for AWS and `gcloud instances create` for GCP 
   can be passed as executor arguments. 

5. Default build image will be either `ubuntu:latest` or `alpine:latest`. If no image specified in the `.gitlab.ci.yml` file, 
   build will still be executed within docker container. 


*Rationale: SSH is portable between different cloud providers*

## MVP Scope

1. AWS implementation only, with code design allowing for extension

*Rationale: GitLab community seems to be having issue with older `docker-machine` implementation for AWS,
so in a way AWS is fix for existing problem, whereas GCP and Azure are extra

2. Implenetation as Custom Gitlab Executor

*Rationale: [Gitlab custom executor](https://docs.gitlab.com/runner/executors/custom.html) interface is well documented and easy to integrate with.
Developing this as plugin would require deeper research and understanding of GitLab runner internals and code. 

3. AWS will support Amazon Linux 2 AMIs only, through default provisioning script working on this OS. Code structure will 
   allow for different OSs, and MVP will include OS discovery. 

4. Cache restore and cache storage is performed from build runners, rather from build manager. Code structure will allow
   for future implementation of cache storage and retrieval from build manager, based on command line argument. 

*Rationale: Restoring and saving distributed cache from manager requires secure file copy operations and introduces
complexity in MVP. Also, large file SCP may not be suitable for large artifacts. *




