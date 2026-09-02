# design_doc_template

A generic design document template for documenting Micro services.

# Title

##### Author: LFSCamargo

## Introduction

_What is the goal of the service? What is the problem statement?_

If there are any product specification documents, link them to appendix.

## Potential Solutions

_Section should include alternative implementations/solutions_

_Document all the Spikes done to evaluate potential solutions. _

_Should start documenting what was done prior to the spike, what are the specific problems and what kind of solution do we value. _

_Is it feasible? If the time runs out, “we don’t know yet” is a perfectly acceptable answer to “is it feasible?”. Document what you tried or learned so any other team member could continue if it is decided to resume the spike. _

_How much effort does it need for each approach? Pro’s and con’s of each approach. The goal of a spike would be to produce knowledge, document results and to reduce uncertainty. _

_Document what were your alternatives, why you made the decision and how that will affect (positively and potentially negatively) the team and the project. _

_Knowledge produced and documented as a result of Spike should be enough for the product owner to create the stories to do the real work and help the development team to estimate it._

## Assumptions

_Clearly document any assumptions you are making to define clear scope of the deliverable._

## Constraints / Limitations

_List of all the things that this system is NOT intended to do or can not do._

## System Design and Architecture

### System diagram or flowchart

_Interaction diagram of various inputs, outputs, sub systems and dependencies._

### Terminology and components

_Nobody likes abbriviations, list of all the key terms and abbriviations_

### Hard and soft dependencies

_What will be the impact on the system if these dependencies go down?_

### Algorithm or Pseudo code for main components

_Describe your logic in this section_

### Service guarantees that you expect

_ACID requirements, Fault tolerance and recovery. If something goes wrong, will the service be able to auto recover or the recovery process has to be manual._

### **Data definition, schema design and Persistence requirements**

_Describe your data model and entity relationships_

### Caching requirements

_What are you planning to cache? What will be size of cache? For how long?_

### Capacity planning

_how much data will service generate, retention policy, sharding_

### Performance requirements

_expected throughput and SLA. How you came up with these numbers?_

### Security

_Any PII / PCI data being handled by service?_

### Multi region story

_Any special requirements your service needs to achieve HA must be listed here. E.g. By default, BigTable doesn’t replicate data across zones or regions. If you are using that to store data, explain how will you handle zone outages.._

## API / gRPC Endpoints

[Instructions](https://github.com/pseudomuto/protoc-gen-doc#writing-documentation) for documenting gRPC methods and services

## Rollout Plan

- Define the roll out phases
- A/B tests that you plan to do as part of rollout
- Data migration

## Test Plan

_How are you planning to test the service? How will you test the dependencies? Will you be creating the mocks for dependencies or using a sandbox environment? How are you planning to run load and performance tests ?_

## Appendix

- References, links to additional documentation (E.g. Product specification or Business requirements with acceptance criteria )

## Review Sign-Off

### Internal Review

- **Team lead/manager**: _Sign-off by _
- **Sponsor**: _Sign-off by _

**Target sign-off:** no late than _2 weeks from date of submission - format MM/DD/YYYY (WWW)._

\*\*
\*\*

**[Sign-off Completed on] **