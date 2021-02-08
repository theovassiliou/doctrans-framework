# Instance ID

This package supports the handling of instance-identification fields as proposed in Theos paper.

In short, via an instance-identification header field, a micro-service instance can disclose it's identity, by responding with structured information.

This becomes in particular of interest, if the micro-service also includes instance-information of other called services in its response.

An single instance is represented through it MIID

`MIID := <sN> "/" <vN> ["/" <vA>] "%" <t>s`

Example: `msA/1.1/feature-branch-2345abcd%222s`

The complete call-graph including it's own MIID is represented by:

```text
CIID := MIID [ "(" UIDs+ ")"]
UIDs := CIID [ "+" CIID ]+
```

This package provides some helpers to work with this type of instance-identification

```text
CIID := MIID [ "(" UIDs+ ")"]
UIDs := CIID [ "+" CIID ]+
MIID := <sN> "/" <vN> ["/" <vA>] "%" <t>s
```

## Supported functionality

This packages supports the

- creation
- parsing
- encoding
- human-friendly display

of instance-id
