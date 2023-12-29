# SimpleAgent

A simple agents system that's used to monitor a list of IP Ranges.

Still playing with some stuff

## What it should do
A Share is a list of ip ranges that the agent scans to determine which ip ranges are being scanned and which are not.
An Agent requests a share based on the machine's capabilities to make sure ip ranges are being scanned on daily basis.

```mermaid
sequenceDiagram
    participant Agent
    participant API
    participant Database
    Agent->>API: Requesting $Count of IP Ranges (Share)
    API->>Database: Retrieve $Count of IP Ranges and mark them
    Database-->>API: Return list
    API-->>Agent: Return Share
    loop IPRanges
        Agent->>Agent: Scan IP Range
    end
    Agent->>API: Save results
    API->>Database: Add results to database
    loop HeartBeat
        API-->>Agent: Send HeartBeat to agents to check for free shares
        API-->>Database: Mark shares as free
    end
```
