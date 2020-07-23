---
id: vpn_config_from_api
title: Configurable VPN from Orchestrator API
hide_title: yes
---

# Configurable VPN from Orchestrator API

- Feature owners: `@alexrod`, `@tdzik`
- Feedback requested from: `@apad`, `@xjtian`, `@shasan`

## Summary

Currently, the only option to gain remote access to any Access Gateway involves some manual steps to be able to introduce an OpenVPN connection to it. This proposal involves some updates to the process to allow a configurable VPN connection through the Orchestrator REST API. 

## Motivation

The current method of configuring and then creating a VPN connection to an AGW, involves some manual steps, detailed:

1. Download and install openvpn package, acting as the OpenVPN client from the AGW. (Which ultimately connects to an AWS instance - OpenVPN server)
2. Register magma AGW on Orchestrator using `show_gateway_info.py` information. (AGW needs to be bootstrapped)
3. Run a script `vpn_setup_cli` which contains fixed server endpoint and port, retrieves the certs from cloud and sets up client config.
4. End-user then uses updated files from setup script to connect to VPN.

We can automate some of this setup by configuring connection parameters from Orchestrator REST API, we also get more control of the setup in case of setup / connection goes wrong by enabling / disabling VPN, or changing the config parameters. Security is also another important topic to take account, by keeping the credentials updated with the gateway certificates / certifier and also revoked once they're not valid or expired.

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
                     | module of Terraform  |                    |                    |   
                     |                      |                    |                    |   
                     |                      |                    +--------------------+   
                     +----------------------+                                          
```

For setting this up, we can take advantage of our terraform module configuration, to deploy and expose an OpenVPN server. 
We can deploy an OpenVPN server on kubernetes by using helm openvpn module on: https://hub.helm.sh/charts/stable/openvpn. This server should use persistent volume in k8s to store all the client keys mapping information. 

VPN credentials will be handed and provisioned from cloud controller during bootstrapping process (which is part of the magmad connection with the access gateway), these should be rotated along with the gateway certs and revoked if these become not valid or are expired. The boostrapper module in cloud will be in charge of managing the client keys stored on disk by a RPC interface, and then hand them down to the gateway only after a successful bootstrap process. We can mount the same persistent volume that the server uses for storing the client keys onto the bootstrapper, which will be the communication and syncing between the bootstrapper and OpenVPN.

From Orchestrator, we can add controller app endpoints that will allow user to do multiple operations on the VPN connection config:
- `.../gateways/gateway_id/vpn_config`
  - port
  - server domain name
  - is enabled

From here, the AGW can spin off and enable an OpenVPN TCP client, we can wrap the client into a dynamic service that can be easier to manage and activate using magmad service. This implementation should give us more flexibility as deploying the OpenVPN server becomes a specification as a Terraform module, provides more security by maintaining the VPN credentials valid along with the Access Gateway certificates, and also provides an scenario for the user to configure and manage the VPN connections right from the API while also help with recovery options when issues arise.

## Timeline of Work

- Adding deployment of OpenVPN server through Terraform module
- Implement interface for bootstrapper management of persistent client kyes
- Update bootstrapper process to include provision / maintenance of VPN creds for enabled VPN gateways
- Add OpenVPN client wrapper to AGWs
- Add cloud endpoints for VPN configuration / management 
