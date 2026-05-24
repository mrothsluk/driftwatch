# driftwatch

Detects infrastructure config drift between Terraform state and live cloud resources and reports diffs.

---

## Installation

```bash
go install github.com/yourorg/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/driftwatch.git && cd driftwatch && go build ./...
```

---

## Usage

Point `driftwatch` at your Terraform state file and let it compare against your live cloud environment:

```bash
driftwatch scan --state terraform.tfstate --provider aws --region us-east-1
```

Example output:

```
[DRIFT] aws_security_group.web (sg-0abc1234)
  - ingress.0.cidr_blocks: ["10.0.0.0/8"] → ["0.0.0.0/0"]

[DRIFT] aws_instance.api (i-0def5678)
  - instance_type: "t3.medium" → "t3.large"

[OK]    aws_s3_bucket.assets
Summary: 2 drifted, 1 clean
```

### Flags

| Flag | Description |
|------|-------------|
| `--state` | Path to Terraform state file |
| `--provider` | Cloud provider (`aws`, `gcp`, `azure`) |
| `--region` | Target cloud region |
| `--output` | Output format: `text`, `json`, `markdown` |
| `--fail-on-drift` | Exit with code 1 if drift is detected (useful in CI) |

### CI Integration

```bash
driftwatch scan --state terraform.tfstate --provider aws --fail-on-drift
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

[MIT](LICENSE)