# trade-tariff-lambdas-xi-certificate-update-mailer

Scheduled go lambda function to notify HMRC CUPID team when certificates on the EU service have changed.

This is anticipated to help with adjustments to the Appendix 5a content.

```mermaid
sequenceDiagram
    participant Scheduler as Scheduler
    participant Lambda as Lambda
    participant AWS S3 Bucket as AWS S3 Bucket
    participant AWS SES as AWS SES
    participant HMRC as HMRC

    Scheduler->>Lambda: Trigger at 08:00 AM
    Lambda->>AWS S3 Bucket: GET XML file from AWS S3 Bucket
    API-->>Lambda: XML file data
    Lambda->>Lambda: Extract certificate updates from XML
    Lambda->>AWS SES: Compose HTML email with certificate update content
    AWS SES->>HMRC: Send HTML email
```
