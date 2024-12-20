# Secure-User-Retriever
This project retrieves a list of VIP user IDs from the BADSEC server using three implementations: Python, Ruby, and Go. The project includes retry logic, error handling, and JSON-formatted output of the user IDs.

## **Project Structure**

```plaintext
noclist-project/
├── src/                        # Source code directory for all implementations
│   ├── python/                 # Python implementation
│   │   └── noclist.py          # Python script
│   ├── ruby/                   # Ruby implementation
│   │   └── noclist.rb          # Ruby script
│   └── go/                     # Go implementation
│       └── noclist.go          # Go script
├── README.md                   # Documentation (this file)
└── COMMENTS                    # Additional comments or instructions
```

## **What Was Done**

1. **Three Implementations**:
   - A Python script (`noclist.py`) using the `requests` library to interact with the BADSEC API.
   - A Ruby script (`noclist.rb`) using `net/http` and `digest` libraries.
   - A Go program (`noclist.go`) leveraging the Go `http` package and SHA256 hashing.

2. **Core Features**:
   - Retrieve an authentication token from the `/auth` endpoint.
   - Generate a checksum for the `/users` endpoint.
   - Fetch and process the list of user IDs.
   - Handle errors with retries (up to 3 attempts per endpoint).
   - Output a JSON-formatted array of user IDs to stdout.

3. **Error Handling**:
   - Retries on transient network failures or non-200 HTTP responses.
   - Exit with non-zero status code on repeated failures.


## **Prerequisites**

1. **Docker**:
   - Ensure Docker is installed and running on your machine. ([Install Docker](https://www.docker.com/get-started))

2. **Programming Languages**:
   - **Python**: Python 3.7+ with `requests` installed.
   - **Ruby**: Ruby installed with default libraries.
   - **Go**: Go 1.16+ installed.

3. **Tools**:
   - Access to a Unix-like command-line environment (Windows PowerShell, Linux Terminal, macOS Terminal).

## **How to Validate**

### **Step 1: Run the BADSEC Server**
1. Start the server using Docker:
   ```bash
   docker run --rm -p 8888:8888 adhocteam/noclist
   ```
2. Confirm the server is running by visiting [http://localhost:8888](http://localhost:8888). You should see:
   ```plaintext
   Listening on http://0.0.0.0:8888
   ```

### **Step 2: Validate Each Implementation**

#### **Python**
1. Navigate to the Python directory:
   ```bash
   cd src/python
   ```
2. Run the Python script:
   ```bash
   python noclist.py
   ```
3. **Expected Output**:
   - A JSON-formatted array of user IDs will appear in the terminal, such as:
     ```json
     ["9757263792576857988", "7789651288773276582", "16283886502782682407"]
     ```

#### **Ruby**
1. Navigate to the Ruby directory:
   ```bash
   cd src/ruby
   ```
2. Run the Ruby script:
   ```bash
   ruby noclist.rb
   ```
3. **Expected Output**:
   - A JSON-formatted array of user IDs will appear in the terminal:
     ```json
     ["9757263792576857988", "7789651288773276582", "16283886502782682407"]
     ```

#### **Go**
1. Navigate to the Go directory:
   ```bash
   cd src/go
   ```
2. Build the Go binary:
   ```bash
   go build noclist.go
   ```
3. Run the binary:
   ```bash
   ./noclist
   ```
4. **Expected Output**:
   - A JSON-formatted array of user IDs will appear in the terminal:
     ```json
     ["9757263792576857988", "7789651288773276582", "16283886502782682407"]
     ```

## **Testing**

### **Manual Tests**
1. **Server Availability**:
   - Test `/auth` by opening `http://localhost:8888/auth` in a browser or using `curl`:
     ```bash
     curl http://localhost:8888/auth
     ```
   - Test `/users` with the checksum header using `curl`:
     ```bash
     curl -H "X-Request-Checksum: <calculated_checksum>" http://localhost:8888/users
     ```

2. **Simulate Failures**:
   - Stop the Docker server and rerun the scripts to test retry logic.


## **Troubleshooting**

1. **Error: Command Not Found (Ruby or Go)**:
   - Ensure Ruby or Go is installed and added to your PATH.
   - Verify installations:
     ```bash
     ruby --version
     go version
     ```

2. **Error: Failed to Fetch Auth Token**:
   - Ensure the Docker server is running:
     ```bash
     docker ps
     ```

3. **Error: Checksum Mismatch**:
   - Verify that the checksum calculation is correct in the scripts.
 
