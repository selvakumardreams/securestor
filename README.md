**SecureStor** is an open source object storage solution that combies high-perfoamnce data storage with integrated **Software Composition Analysis (SCA)** to ensure security and compilance. It provided advanced features like **AI-driven data tiering, end-to-end encryption, automated compilance auditing, version control for objects and hybrid cloud replication**, empowering businesses to store, manage, and protect their data while minimizing risks from software vulnerabilities and optimizing costs.

## Todo
1. **Built-in Software Composition Analysis (SCA)** - This feature would automatically scan and track vulnerabilities in all dependencies used by the system and alert users in real time.
2. **Automated Data Compliance Auditing** - Compliance-as-a-Service: Many industries require strict compliance with data regulations (GDPR, HIPAA, etc.). It could offer built-in compliance auditing tools that automatically check and enforce compliance policies on stored data, providing reports to users on how their data storage conforms to relevant regulations.
3. **Hybrid and Multi-cloud Data Replication** -  Users can configure their object storage to automatically replicate data across multiple cloud providers or across on-premises and cloud environments.
4. **AI/ML-Driven Data Tiering & Optimization** - Use AI/ML algorithms to analyze the usage patterns of data and automatically categorize objects based on frequency of access. The system could then tier data intelligently (hot, cold, and archive) to optimize storage costs and performance.
5. **End-to-End Encryption with Key Management** - Users have complete control over their encryption keys. Implement a bring-your-own-key (BYOK) system where users can integrate their own hardware security modules (HSMs) or key management services (KMS) to further enhance security.
6. **Advanced Search and Metadata Indexing** - Ability to search objects in storage based on metadata, tags, and custom indexes
7. **Version Control for Objects** - Offer version control for objects, enabling users to store multiple versions of the same object and track changes over time.
8. **Serverless Object Storage Integration** - Support serverless functions directly within the object storage platform. Users could define serverless workflows that trigger based on certain events—like uploading a file, deleting an object, or reaching a storage threshold.
9. **Data Provenance & Integrity Verification** - Offer data provenance tracking, allowing users to trace the origin and history of any stored object. This could include metadata on how and when the object was created, modified, and accessed. Integrity verification mechanisms, such as block-level checksums or Merkle trees,
10. **Customizable Storage Policies** - Provide highly customizable storage policies that allow users to define retention policies, access rules, and cost optimization strategies based on user, group, or object type.
11. **Data Lifecycle Automation** - Go beyond simple object retention by offering data lifecycle automation that automatically moves data to the appropriate storage tiers or archives based on usage, business rules, or compliance requirements.
12. **Immutable Storage for Compliance and Security** - Add the ability to create immutable storage buckets where objects cannot be deleted or modified for a predefined retention period. This would be particularly useful for industries that require immutable backups for compliance with regulations like financial audits or healthcare data retention.

## Instructions
1. Run the Server: Start the server by running the following command in your project directory.

```go
cd cmd\server
go run main.go
```

```npm
cd frontend
npm run start-legacy
```

2. Create a Bucket: Use curl or a similar tool to create a bucket.

```
curl -X POST "http://localhost:8080/create-bucket?bucket=mybucket"
```

3. Upload a File: Use curl or a similar tool to upload a file to the bucket.

```
curl -F "file=@/path/to/your/file" "http://localhost:8080/upload?bucket=mybucket"
```

4. Download a File: Use curl or a web browser to download a file from the bucket.

```
curl -O "http://localhost:8080/download?bucket=mybucket&filename=yourfile"
```

5. List Files in a Bucket: Use curl or a web browser to list all files in the bucket.

```
curl "http://localhost:8080/list?bucket=mybucket"
```

6. To add / update / delete custom metadata, send a POST request to /update-metadata with the following JSON body:

```
{
    "id": "file-id",
    "custom_metadata": {
        "newKey": "newValue"
    },
    "action": "add"
}
```
```
{
    "id": "file-id",
    "custom_metadata": {
        "keyToDelete": ""
    },
    "action": "delete"
}
```
```
curl "http://localhost:8080/update-metadata"
```