# Security Policy

## Reporting a Vulnerability

We take the security of SQLens seriously. If you believe you have found a security vulnerability, please report it to us by following these steps:

1. **Do not open a GitHub issue.** 
2. Email your report to: `security@example.com` (Placeholder)
3. Include as much detail as possible: steps to reproduce, impact, and a proof of concept if available.

We will acknowledge your report within 48 hours and keep you updated on our progress.

## Redact Mode

SQLens includes a `RedactSensitive` mode (set via `SQLENS_REDACT_SENSITIVE=true`). When enabled, the proxy will mask all literals in SQL queries before they reach the dashboard or logs. We recommend enabling this in production environments where queries may contain PII.
