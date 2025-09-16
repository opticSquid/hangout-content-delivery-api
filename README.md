# HANGOUT CONTENT DELIVERY API

This service provides temporary time bound access to various content in the platform. This service is tightly integrated to AWS.

- Uses presigned URL to give temporary access to images
- Uses presigned cookies to give temporary access to video chunks

## Performance Requirements

This service is called in repeated succession for all the contents in the platform that user is and will be seeing. So, it is paramount for this service to be extremely fast. Its P90 latency should less than 500ms.

## Integrations

 This service is integrated with standard observability pipeline using open telemetry