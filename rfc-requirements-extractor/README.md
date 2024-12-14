rfc-requirements-extractor
================

A simple tool to convert an XML of RFC to a tree of directories broken by RFC section containing text files of the requirements, one file per mandate.

## Usage

To generate the file tree for RFC 9110, run the following command:

```sh
go run ./main.go -input ../rfcs/rfc9110.xml -output-dir ../rfcs/9110/text
```

## TODO:

- [ ] Parse `xref` tags and convert them into links
- [ ] Bug: `ul` following `t` are missed because the it's hard to know the `ul` is relevant to the `t`

# Helpful References

- [RFC XML Vocabulary](https://authors.ietf.org/rfcxml-vocabulary)
