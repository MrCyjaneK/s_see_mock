# mock format

## v2.0.0

```json
{
    "sha1":   hex.EncodeToString(sha1H.Sum(nil)),
    "sha224": hex.EncodeToString(sha224H.Sum(nil)),
    "sha256": hex.EncodeToString(sha256H.Sum(nil)),
    "sha512": hex.EncodeToString(sha512H.Sum(nil)),
    "md5":    hex.EncodeToString(md5H.Sum(nil)),
}
```