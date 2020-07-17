---
id: 202007_vpn_config_from_api
title: Configurable VPN from Orchestrator API
status: in_review
hide_title: yes
---

# Configurable VPN from Orchestrator API

## Summary

The proposal of creating a central configuration component within Thanos has been rejected by the community as the requirements are specific to the use case at Improbable and that adding configuration management into Thanos will result in adding more knowledge to the system about what the scrapers are doing and their targets.

Cluster configuration for targets will be implemented in a separate repository and we may look at open sourcing in the future if there are others that have the same needs.

## Motivation

Currently, each scraper manages their own configuration via [Prometheus Configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) which contains information about the [scrape_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cscrape_config%3E) and targets that the scraper will be collecting metrics from.

As we start to dynamically scale the collection of metrics (new Prometheus instances) or increase the targets for a current tenant we wish to keep the collection of metrics to a consistent node and not re-allocate shards to other scrapers.

## Goals

The primary goal of having a specific `thanos config` component are:

- enable (auto)scaling of scrapers to meet the needs metrics ingestion
  - ensure we can schedule work to new scrapers
  - on scale down ensure we can drain scrapers and upload data before removing
- assign targets to scrapers based on utilisation

## Proposal

As we scale out horizontally we wish to manage configuration such as target assignment to scrapes in a more intelligent way, firstly by keeping targets on the same scrape instance as much as possible as we scale the scrape pool and in the future more intelligent bin packing of targets onto scrapers.

We would look to add a new component to the Thanos system `thanos config` that would be a central point for configuring what each Prometheus scrapes for the entire cluster and advertising via APIs the configuration for each sidecar.

The `thanos config` will have a configuration endpoint that each `thanos sidecar` will call into to get their own scrape_config jobs along with their targets. Once the sidecar has its jobs it will be able to update targets / scrape_config for the Prometheus instance it is running alongside. This update will primarily be based on `file_sd_config` and will be allowed to add or remove targets without changing Prometheus itself.

The config component will keep track of what sidecar's are in the cluster via the existing gossip mechanism and will therefore have a central view of targets to Prometheus instances.
When we scale up our scrape pool and a new scrape instance comes online the `thanos sidecar` join the gossip cluster and therefoe `thanos config` will know that it can start assigning configuration to this new node.
In a scale down scenario we can first remove all targets from a given scraper effectively draining that instance of work and kick off the process of uploading the WAL to storage (not in scope of this design). During the time that the Prometheus instance has no targets we would still want to be able to query the instance for data that has not yet been uploaded.

We believe that a central point for configuration and management is better in this scenario as it gives us more flexibility in the future to add bin packing / consistent hashing. It would also be an ideal place for deciding on "hot shard" issues, the config component would be able to see the utilization of each node and decide based on that where to schedule work. Having a centralised approach would also help with debugging, testing and maintaining the code when issues arise.


## Timeline of Work
