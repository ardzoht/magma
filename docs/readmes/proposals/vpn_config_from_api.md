---
id: vpn_config_from_api
title: Configurable VPN from Orchestrator API
status: in_review
hide_title: yes
---

# Configurable VPN from Orchestrator API

## Summary

- Currently, the only option to gain remote access to any Access Gateway involves some manual steps to be able to introduce an OpenVPN connection to it. This proposal involves some updates to the process to allow a configurable VPN connection through the Orchestrator REST API. 

## Motivation

- The current method of configuring and then creating a VPN connection to an AGW, involves some manual steps, detailed:

1. Download and install openvpn package, which ends up being the OpenVPN client from the AGW. (Which ultimately connects to an AWS instance - OpenVPN server)
2. Register magma AGW on Orchestrator using `show_gateway_info.py` information. (AGW needs to be bootstrapped)
3. Run a script `vpn_setup_cli` which contains fixed server endpoint and port, retrieves the certs from cloud and sets up client config.
4. End-user then uses updated files from setup script to connect to VPN.

- We can automate some of this setup by configuring connection parameters from Orchestrator REST API, we also get more control of the setup in case of setup / connection goes wrong by enabling / disabling VPN, or changing the config parameters.

- Security is also another important topic to take account, by keeping the credentials updated with the gateway certificates / certifier and also revoked once they're not valid or expired.

## Goals

The goals of automating VPN setup and making it configurable are:

- Allow an easier method of configuring VPN connection setup for AGWs
- Automating enabling / disabling of the connection
- Improve security of the VPN workflow

## Proposal

```                                                                                         
|-------------------------------|                        +-------------------------------+
|                               |      VPN Credentials   |                               |
|                               |------------------------|                               |
|                               |    magmad connection   |         Access Gateway        |
|   Orchestrator Bootstrapper   |                        |                               |
|                               |                        |                               |
|                               |                        |                               |
|                               |                        |                               |
|                               |                        |                               |
|--------------------------------                        +-------------------------------+
                                                                          |               
                                                                          |               
                                                                          /               
                                                                         |                
                                                                         |                
                     +----------------------+                    +--------------------+   
                     |                      |                    |                    |   
                     |                      |                    |                    |   
                     |    OpenVPN Server    |--------------------|  TCP OpenVPN Client|   
                     |                      |                    |                    |   
                     |                      |                    |                    |   
                     |                      |                    +--------------------+   
                     +----------------------+                                          
```

- For setting this up, we can take advantage of our terraform module configuration, to spin off and deploy an OpenVPN server. Feedback required from @xjtian on how to add this specification.

- From Orchestrator, we can add endpoints 

As we scale out horizontally we wish to manage configuration such as target assignment to scrapes in a more intelligent way, firstly by keeping targets on the same scrape instance as much as possible as we scale the scrape pool and in the future more intelligent bin packing of targets onto scrapers.

We would look to add a new component to the Thanos system `thanos config` that would be a central point for configuring what each Prometheus scrapes for the entire cluster and advertising via APIs the configuration for each sidecar.

The `thanos config` will have a configuration endpoint that each `thanos sidecar` will call into to get their own scrape_config jobs along with their targets. Once the sidecar has its jobs it will be able to update targets / scrape_config for the Prometheus instance it is running alongside. This update will primarily be based on `file_sd_config` and will be allowed to add or remove targets without changing Prometheus itself.

The config component will keep track of what sidecar's are in the cluster via the existing gossip mechanism and will therefore have a central view of targets to Prometheus instances.
When we scale up our scrape pool and a new scrape instance comes online the `thanos sidecar` join the gossip cluster and therefoe `thanos config` will know that it can start assigning configuration to this new node.
In a scale down scenario we can first remove all targets from a given scraper effectively draining that instance of work and kick off the process of uploading the WAL to storage (not in scope of this design). During the time that the Prometheus instance has no targets we would still want to be able to query the instance for data that has not yet been uploaded.

We believe that a central point for configuration and management is better in this scenario as it gives us more flexibility in the future to add bin packing / consistent hashing. It would also be an ideal place for deciding on "hot shard" issues, the config component would be able to see the utilization of each node and decide based on that where to schedule work. Having a centralised approach would also help with debugging, testing and maintaining the code when issues arise.


## Timeline of Work
