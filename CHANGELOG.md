# Changelog

## Unreleased

### Compatibility

- SigV2 now signs bucket-only paths exactly as sent: `/bucket` and `/bucket/`
  are distinct. Clients that send `/bucket` but sign `/bucket/` must update
  their signing logic or send the trailing slash. The verifier no longer adds
  a trailing slash implicitly.
