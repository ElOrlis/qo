# 📝 qo Package TODO List

## 🎯 Phase 1: Core Enhancements
- [ ] Add **custom retry policies** (linear, exponential, constant backoff)
- [ ] Allow **configurable retryable status codes**
- [ ] Implement **max retry attempts** as a user-defined setting

## 📜 Phase 2: Logging & Debugging
- [ ] Integrate structured logging (`zap` or `logrus`)
- [ ] Enable **request/response logging** for debugging
- [ ] Implement a **verbose/debug mode** to log retry attempts and errors

## 🔧 Phase 3: Request Enhancements
- [ ] Support **custom headers and middleware**
- [ ] Add **timeout handling** for requests
- [ ] Modify `do` method to accept **request bodies** for `POST`, `PUT`, `PATCH`
- [ ] Implement **GZIP compression for request bodies**
- [ ] Support **automatic decompression of GZIP responses**
- [ ] Allow **custom compression algorithms** (e.g., Brotli, Deflate)

## ⚡ Phase 4: Rate Limiting & Circuit Breaking
- [ ] Implement a **rate limiter** to prevent excessive requests
- [ ] Integrate **circuit breaker pattern** (`sony/gobreaker`)

## 💾 Phase 5: Caching Mechanism
- [ ] Implement **in-memory caching** for frequent requests
- [ ] Allow **custom cache backends** (Redis, file-based)

## 🔄 Phase 6: Automatic Retries & Hooks
- [ ] Differentiate retry strategies for **idempotent vs. non-idempotent methods**
- [ ] Add **before/after request hooks**
- [ ] Implement **retry event hooks** for logging or monitoring

## 🔐 Phase 7: Security & Failover
- [ ] Implement **TLS support** with configurable settings (TLS versions, certificates)
- [ ] Add **automatic failover logic** for multiple endpoints:
  - [ ] Allow users to define **multiple fallback URLs**
  - [ ] If primary fails, retry on secondary and tertiary endpoints

## 🔑 Phase 8: Authentication Support
- [ ] Provide **OAuth2 token support**
- [ ] Add **API key authentication support**

## 🔨 Phase 9: CLI Tool
- [ ] Build a simple **CLI tool** to test requests and retry behavior

## 📖 Phase 10: Documentation & Examples
- [ ] Write a **detailed README** with usage examples
- [ ] Generate **GoDoc documentation**
- [ ] Create **example implementations** and real-world use cases

