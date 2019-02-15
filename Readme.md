# K8S-job-sidecar-terminator
It's workaround for specific condition when side car mode cloudsql with K8S job  
the project is for [This Issue](https://github.com/kubernetes/kubernetes/issues/25908)
**And do not expose the default port 8080**  

## Usage:
In cloudsql-proxy:  
`job-terminator {your command}`  
and kill this container in other job in the same pod  
`{your command} && curl 127.0.0.1:8080/kill`  
or build your own images and:  
`job-terminator --killer {your command}`  
