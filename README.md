# IRA Service

Introducing Auto-Remediation and Other Automation into 
Chainlink's Incident Response Process

The goal of this RFC is to evaluate various solutions for automating incident response. One solution is chosen and explored in depth, but each explored solution has pros and cons that need to be considered.

**Authors**: @KA-ROM (Kim Romero)

## Motivation

Chainlink does not have an incident auto-remediation service for production incidents. This builds tech debt since our fully manual incident response process decreases stakeholders’ incentive for building self-healing alerting. For example, many triaging steps could be automated, but currently are not. Additionally, only one Chainlink product is close to a design that allows auto-remediation (VRF).

The result of this is a slower Time to Event (TTE) for all incidents due to the manual triaging. Additionally, engineers are paged into war rooms sometimes fixed by processes that could also be automated. 

However, there is no service for this logic to fit into appropriately.



## Goals

* Agreeing upon the best solution for introducing alert auto-remediation into Chainlink Lab’s existing alerting stack
* Increase reliability by introducing self-healing design into our incident response process
* Reduce TTE by adding more intelligent automation and information-amending to alerts received by human responders

## Anti goals

This RFC focuses on creating the MVP of this service. To achieve this, the MVP will be defined to contain the following two goals:

1. Full alert auto-remediation: one method that auto fires on <get alert name> and runs the automation of the full VRF runbook. On failure, the alert still escalates to a human responder with the remediation progress amended to the original alert (or link to the service’s output UI).
2. Partial triaging: amendment of all feed alerts with the data included in IRA's feed report command (an IRA command that combines crucial feed configuration data).

This RFC is not concerned with all the future automation this service may include.




## Glossary
**IR:** The Incident REsponse team.

**IRA:** Existing golang CLI created by @KA-ROM and contributed to by the IR team that houses triaging automations, as well as other automations.

**PD:** PagerDuty



## Potential Solutions

Qualities a Solution Must Have:

* **Customizability automation:** The solution should be highly automated to minimize manual intervention and reduce response time. To start, it must enable the following:
* **Full alert auto-remediation:** Enable IR to automate all triaging steps in alerts when possible.
* **Partial alert automation:** Enable IR to automate some triaging steps so that PagerDuty alerts are immediately amended with necessary data.
* **Scalability:** The solution should be able to handle large-scale incidents and work seamlessly with other tools in the company's infrastructure.
* **Resilience:** The solution should be designed to be resilient to prevent the service from being disrupted in case of incidents.
* **Security:** The solution should comply with relevant security standards and regulations to prevent unauthorized access and data breaches.
* **Monitoring:** The solution should have built-in monitoring capabilities to track the performance and effectiveness of the auto-remediation service.

### Solutions

#### On-Premise Solution: 
A self-managed auto-remediation service that runs in a Kubernetes cluster within Chainlink Lab’s infrastructure. The service would have a web server that can receive alerts from Alert Manager and trigger the appropriate remediation functions based on the alert type.

Key considerations for the architecture of such a service:

**Event-driven architecture:** this should be designed with an event-driven architecture that can handle multiple types of alerts from different sources. This can be achieved by implementing a message broker such as Apache Kafka or RabbitMQ to handle incoming alerts and distribute them to the appropriate remediation function.

**Scalability:** the service should be designed to scale horizontally as the number of alerts and remediation functions grow. 

**Security:** should be designed with security in mind, including secure communication channels, encryption of sensitive data, and access control mechanisms to prevent unauthorized access to the service.

**Monitoring:** the service must be equipped with monitoring and logging capabilities to track the performance and effectiveness of the remediation functions.

**Pros:**
* Greater control: full control over the design and functionality of the service allows us to tailor it to our specific needs and requirements. Remediation functions can be customized to our specific environment and use cases, potentially leading to more effective and efficient remediation.
* Cost-effective: Using in-house tools can be more cost-effective in the long run than paying for third-party solutions.

**Cons:**
* Maintenance: we are responsible for maintaining the service, including updating dependencies, monitoring, and responding to any issues.
* Complexity: building and maintaining an in-house service can be more complex than using a managed solution, and may require additional resources and expertise.
* Full responsibility for errors: developing and managing an in-house service increases the potential for human error, which could lead to issues such as misconfigured remediation functions or security vulnerabilities.

---


#### AWS Systems Manager:
Automation through this AWS service enables the use of pre-written and custom runbooks. AWS provides ready-made operational runbooks (called automation documents) that cover many predefined use cases. In cases where there are no ready-made automation documents, the solution includes additional AWS Lambda functions that require minimal code and no external dependencies outside AWS native tools.

The incident response actions are executed in a central security account, so there is no need to change the service accounts where incidents are monitored and responded to. This solution can also include exception handling when actions should not be executed, and manual triggering of remediations that should not receive an automatic response. 

**Pros:**
* Offers a variety of pre-built services for automation.
* Highly scalable.
* Highly flexible.

**Cons:**
* Additional cost for using the services.
* Integration with existing tools may be challenging.
* Requires additional training and expertise to set up and manage the AWS infrastructure.
* Risks micro-service madness. 



**AWS Lambda (only):**
AWS Lambda is a serverless computing service that enables developers to run code without managing servers. This solution requires writing custom code for each incident response scenario.

**Pros:**
* Complete control over the incident response process with the ability to customize code for each scenario.
* No need to rely on third-party services or tools.
* Cost-effective as AWS Lambda is charged per execution time.

**Cons:**
* Requires writing custom code for each incident response scenario, which can be time-consuming.
* Needs careful testing and ongoing maintenance to ensure that the code functions correctly and does not introduce new problems.
* Limited visibility and traceability into the incident response process as there is no central dashboard or tool for managing the process.
* May become expensive on high alert volumes.
* Lack of clear central oversight, potentially reducing coordination. 

---

#### PagerDuty Playbooks: 
A PagerDuty that enables teams to create predefined workflows for specific incidents. Each playbook consists of a series of steps that guide users through the resolution process. Playbooks can include multiple paths based on the incident type, severity, and other factors. They can also incorporate automated actions, such as triggering an AWS Lambda function or sending an email notification.

Regarding implementing blockchain queries and other similar needs, PagerDuty playbooks can integrate with various third-party tools and services through its integration ecosystem.

**Pros:**
* Already heavily used. Playbooks provide a centralized and standardized way of responding to incidents.
* Integration with PagerDuty is straightforward.
* Offers pre-built playbooks for automating alert triage.

**Cons:**
* Potentially limited customization options.
* Likely to need additional customization or integration with other tools to involve blockchain queries
* The cost of PagerDuty's services may be prohibitive 

---

#### StackStorm:
An open-source event-driven automation platform that integrates with popular monitoring, alerting, and communication tools to automate and coordinate incident response. It offers a wide range of integrations, offers basic plugins, and supports the execution of arbitrary code.

**Pros:**
* Provides a wide range of integrations with popular tools such as Prometheus, Alert Manager, PagerDuty, and Grafana
* Supports the execution of arbitrary code
* Built on an event-driven architecture, which means it can react to events in real time, offering event-driven automation
**Cons:**
* May have a steep learning curve for team members who are not familiar with the tool
* Requires a dedicated server to run the StackStorm service


---



#### SaaS:
I considered TFTT, Zapier, and < >, but we would be trusting them as third-party providers with sensitive data and code that future automation would benefit from including. You would not have full control over the code’s access or its data, leaving us at risk of a security breach. On top of the cost of a SaaS solution, this section will not be expanded upon. 

---



## Risks and Failure Scenarios

**Service disruptions:** 
The auto-remediation service may experience service disruptions due to infrastructure issues, such as network outages or hardware failures. This could result in delayed incident response times or complete failure of the auto-remediation service, and demand a regression to manual triaging.

To mitigate this risk, the service should be designed to be resilient and fault-tolerant, with backup systems and contingency planning in place. 

**Software bugs:** The auto-remediation service may contain software bugs that cause unexpected behavior or failures. This could lead to incorrect incident response actions or failure to detect and respond to incidents.

To mitigate this risk, the service should be thoroughly tested and validated before deployment, and ongoing testing and monitoring should be conducted to detect and address any issues that arise.

**False negatives:** The auto-remediation service may fail to detect or respond to incidents that require intervention, leading to prolonged downtime or other issues. This could result in decreased system availability and increased downtime.

To mitigate this risk, the service should be designed to be flexible and customizable, with the ability to adapt to changing circumstances and respond to new types of incidents as they arise.

**Integration issues:** The auto-remediation service may not integrate seamlessly with other tools or systems in the company's infrastructure, leading to delays or failures in the incident response process. 

To mitigate this risk, the service should be designed to be compatible with existing systems, and thorough testing and validation should be conducted before deployment.



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



