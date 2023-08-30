# trade-tariff-lambdas-xi-certificate-update-mailer

Scheduled go lambda function to notify HMRC CUPID team when certificates on the EU service have changed. This is anticipated to help with adjustments to the Appendix 5a content.

```mermaid
sequenceDiagram
    participant Scheduler as Scheduler
    participant Lambda as Lambda
    participant API as API
    participant SES as SES
    participant HMRC as HMRC

    Scheduler->>Lambda: Trigger at 08:00 AM
    Lambda->>API: GET XML file from https://webservices.hmrc.gov.uk/taric/2022-05-14_TGB22134.xml
    API-->>Lambda: XML file data
    Lambda->>Lambda: Extract certificate updates from XML
    Lambda->>SES: Compose HTML email with certificate update content
    SES->>HMRC: Send HTML email
```
