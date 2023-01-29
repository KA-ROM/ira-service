# IRA Service
This RFC focuses on the creation of the MVP of an IR service that introduces auto-remidation, data collection, and other automations into Chainlink's IR process.


**Authors**: @KA-ROM


## Motivation


* Currently, alerting logic is not maintained well. Each product team owns their alerting across different products (chain-monitor, flux-emulator, evm-exporter, alert-manager) and there isn’t an alert review process in place nor method for IR to work recursively with other teams to refine alerts. This has ultimately resulted in alerting tech-debt that affects our reliability.
* Many triaging steps could be automated but currently are not. This means slower TTE for all incidents. Some alerts could be fully resolved programmatically but currently are not, but there currently is no service for this logic to fit appropriately into.
* Data collection is difficult outside of alerts that escalate into war-room situations and reach Jira documentation.


## Goals




The goal of this document is to outline the path forward for the creation of a new IR service that enables alert auto-remediation, data collection, and other automations to be possible in a manner that IR owns for the sake of exploring an introduction of automations and auto-remidations into Chainlink's IR process. A service would enable IR to develop:


* Full alert auto-remediation: Enable IR to fully automate all triaging steps in the alerts when possible.
* Partial alert automations: Enable IR to automate some triaging steps so that PagerDuty alerts are immediately amended with necessary data.
* Lessen noise: allow IR to create meaningful automations when product teams cannot immediately deliver observed improvements (example: checking RDD config values rather than using hardcoded values).
* Data-collection: Collect fine-grain information on alert data beyond those that escalate into war-rooms and make it into Jira.




### Anti-Goals


This RFC focuses on creating the MVP of this service by focusing on two first goals:
1. Auto-remediation of threshold-deviation alerts that are false-positive due to using the hardcoded value rather than the feed's specific value defined in RDD.
2. Amendment of all feed alerts with the data included in IRA's `feed report` command.
3. TO CONSIDER: data collection of product frequency?


This RFC is not concerned with all the future automations that this service may include.


# Glossary


* **IR**: The Incident REsponse team.
* **IRA**: Existing golang CLI created by @KA-ROM and contributed to by the IR team that houses triaging automations, as well as other automations.
* **PD**: PagerDuty


# Overview of the Current Alerting Stack Across Products


### VRF
### Keepers
### Data Feeds


# Design


### Proposed changes to alert schemas


### Proposed changes to PagerDuty


### Proposed changes to the response process


### Injecting additional data
### Full auto-remediation
### Partial auto-remediation






## Reliability
## Risks and Failure Scenarios
## Alternative Designs
# External Dependencies
# Open Questions



