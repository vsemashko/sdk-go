# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.38.x  | :white_check_mark: |
| 1.37.x  | :white_check_mark: |
| < 1.37  | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to [security@temporal.io](mailto:security@temporal.io).

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

Please include the following information in your report:

- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

This information will help us triage your report more quickly.

## Security Best Practices

When using the Temporal Go SDK, follow these security best practices:

### 1. TLS Configuration

**Always enable TLS in production environments:**

```go
import (
    "crypto/tls"
    "go.temporal.io/sdk/client"
)

// Good: Use TLS with proper certificate verification
opts := client.Options{
    HostPort: "your-temporal-cluster.temporal.io:7233",
    ConnectionOptions: client.ConnectionOptions{
        TLS: &tls.Config{
            // Use system CA pool or custom CA
        },
    },
}
```

**Never disable hostname verification in production:**

```go
// BAD: Do NOT do this in production
conf.InsecureSkipVerify = true  // Makes you vulnerable to MITM attacks
```

### 2. Credential Management

**Use secure credential storage:**

```go
// Good: Use credential providers or secure secret management
opts.Credentials = client.NewAPIKeyDynamicCredentials(func(ctx context.Context) (string, error) {
    // Retrieve from secure secret manager
    return retrieveFromVault(ctx, "temporal-api-key")
})

// Avoid: Environment variables for production credentials
// apiKey := os.Getenv("TEMPORAL_API_KEY")  // Can leak through logs
```

**Protect certificate files:**

```bash
# Ensure proper file permissions for certificate files
chmod 400 /path/to/client.key
chmod 444 /path/to/client.crt
```

### 3. Timeouts and Resource Limits

**Configure appropriate timeouts:**

```go
// Configure workflow execution timeout
opts := workflow.ActivityOptions{
    StartToCloseTimeout: 10 * time.Minute,
    HeartbeatTimeout:    30 * time.Second,
}

// Configure client timeouts
clientOpts := client.Options{
    ConnectionOptions: client.ConnectionOptions{
        // Timeouts are configured automatically, but can be customized
    },
}
```

### 4. Input Validation

**Always validate workflow and activity inputs:**

```go
func MyWorkflow(ctx workflow.Context, input string) error {
    // Validate input
    if len(input) > maxInputSize {
        return fmt.Errorf("input exceeds maximum size")
    }
    if !isValidInput(input) {
        return fmt.Errorf("invalid input format")
    }
    // ... proceed with workflow
}
```

### 5. Error Handling

**Avoid exposing sensitive information in errors:**

```go
// Good: Generic error for external consumption
return fmt.Errorf("authentication failed")

// Avoid: Detailed error that could leak information
// return fmt.Errorf("authentication failed for user %s with key %s", user, apiKey)
```

### 6. Data Encryption

**Use payload encryption for sensitive data:**

```go
import (
    "go.temporal.io/sdk/converter"
)

// Use encryption codec for sensitive workflow data
dc := converter.NewCodecDataConverter(
    converter.GetDefaultDataConverter(),
    NewEncryptionCodec(), // Implement your encryption
)

opts := client.Options{
    DataConverter: dc,
}
```

### 7. Dependency Management

**Keep dependencies updated:**

```bash
# Regularly update dependencies
go get -u go.temporal.io/sdk@latest
go mod tidy

# Check for known vulnerabilities
go list -json -m all | nancy sleuth
```

### 8. Network Security

**Use private networks and firewalls:**

- Deploy Temporal server in a private network
- Use firewall rules to restrict access
- Implement network segmentation
- Use VPNs for remote access

### 9. Audit and Monitoring

**Enable comprehensive logging and monitoring:**

```go
import (
    "log/slog"
    "go.temporal.io/sdk/log"
)

// Use structured logging
logger := log.NewStructuredLogger(
    slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo, // Don't use Debug in production
    })),
)

opts := client.Options{
    Logger: logger,
}
```

**Monitor for suspicious activity:**
- Failed authentication attempts
- Unusual workflow patterns
- Resource exhaustion
- Network anomalies

### 10. Least Privilege

**Run with minimal permissions:**

```bash
# Run worker processes with dedicated service accounts
# Limit permissions to only what's necessary
```

## Security Checklist

Before deploying to production, verify:

- [ ] TLS is enabled and properly configured
- [ ] Hostname verification is enabled (`InsecureSkipVerify = false`)
- [ ] Credentials are stored securely (not in environment variables or code)
- [ ] Certificate files have appropriate permissions (0400 for keys)
- [ ] Timeouts are configured to prevent resource exhaustion
- [ ] Input validation is implemented for all workflows and activities
- [ ] Sensitive data is encrypted using payload encryption
- [ ] Error messages don't expose sensitive information
- [ ] Dependencies are up-to-date and scanned for vulnerabilities
- [ ] Logging level is appropriate (INFO or WARN, not DEBUG)
- [ ] Monitoring and alerting are configured
- [ ] Network access is restricted using firewalls
- [ ] Services run with least privilege

## Additional Resources

- [Temporal Security Documentation](https://docs.temporal.io/security)
- [Go Security Best Practices](https://go.dev/doc/security/best-practices)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)

## Security Updates

Security updates will be released as patch versions. Subscribe to our [GitHub releases](https://github.com/temporalio/sdk-go/releases) or security advisories to stay informed.

## Acknowledgments

We appreciate the security research community's efforts in responsibly disclosing vulnerabilities. Security researchers who report valid security issues will be acknowledged in our security advisories (with permission).
