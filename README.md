Bluenoise is an open source tool for high performance object storage and management of artifacts SBOM and vulnerabilities scanning


## Instructions
1. Run the Server: Start the server by running the following command in your project directory.

```go
go run main.go
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