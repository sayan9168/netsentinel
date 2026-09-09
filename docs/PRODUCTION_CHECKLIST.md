# AOXNet Production Readiness Checklist

## Required before production

- [ ] Run `gofmt -w .`
- [ ] Run `go vet ./...`
- [ ] Run `go test ./... -count=1`
- [ ] Run `go test -race ./... -count=1`
- [ ] Run fuzz tests for frame/control decoders
- [ ] Run benchmarks and inspect allocations
- [ ] Deploy behind TLS 1.3
- [ ] Use mTLS for mutually authenticated peers where required
- [ ] Configure certificate trust and identity validation
- [ ] Configure payload, stream, and connection limits
- [ ] Configure read/write/handshake deadlines
- [ ] Configure keepalive thresholds
- [ ] Monitor protocol errors, bytes, frames, and active streams
- [ ] Load-test concurrent streams and slow consumers
- [ ] Review dependency and Go toolchain updates
- [ ] Perform an application-specific security review

## Compatibility gates

- Existing frame types retain their semantics.
- Mandatory capabilities are negotiated explicitly.
- Unknown optional capabilities are ignored safely.
- Unsupported protocol versions fail deterministically.
- Stream IDs and control frames are validated before state mutation.

## Security gates

CRC32 is not cryptographic authentication. Never expose AOXNet directly to an untrusted network without an authenticated secure transport. TLS certificate validation and application authorization remain deployment responsibilities.
