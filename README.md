
HTTP Compliance Testing
===========================

This repository aims to contain the input/output data comprising the HTTP compliance test suite. It also contains a simplistic script to extract the rules from the XML.

## Background

This project is based on the two conversations, one took place on the HTTP Working Group mailing list and the other on the HTTP Workshop. The blog post [Towards Validated HTTP Implementation](https://www.caffeinatedwonders.com/2024/12/18/towards-validated-http-implementation/) summarizes both discussions and gives a high-level overview of the project.

## RFC Requirement -> Log Level Mapping

Test suite execution engines (TBD: To Be Developed) should map the requirements violations to log levels and results as follows:

| Requirement     | Log Level |
|-----------------|-----------|
| MUST            | ERROR     |
| MUST NOT        | ERROR     |
| REQUIRED        | ERROR     |
| SHALL           | ERROR     |
| SHALL NOT       | ERROR     |
| SHOULD          | WARNING   |
| SHOULD NOT      | WARNING   |
| RECOMMENDED     | WARNING   |
| NOT RECOMMENDED | WARNING   |
| MAY             | INFO      |
| OPTIONAL        | INFO      |

The mapping table is built according to the [RFC 2119](https://datatracker.ietf.org/doc/html/rfc2119).

## Development

The proecss is yet to be refined. Currenty the XML files are placed in the `rfcs` directory. For each RFC, a subdirectory is created under `rfcs`. The subtree has immediate descendants `text` and `data`. The `text` content is generated by the program in `rfc-requirements-extractor`. The `data` sub-directory is meant to contain the codified data corresponding to the requirements as defined in the RFC.

### TODO:

- [ ] Refine the process
- [ ] Define the data format & structure as suitable for the respective RFC section
	- [ ] Define the test cases

### Tip
To find rules in the XML, search using the regular expression:

```regexp
<bcp14>\w+(\s\w+)?</bcp14>
```
