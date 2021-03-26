# URI

```
         _            _  
o  / / _|     URI      |_ 
o / /   |_ Templates  _|
```

> A commandline tool to resolve URI Templates expressions as specified in [RFC 6570](http://tools.ietf.org/html/rfc6570).

For a complete syntax reference check the RFC6570 specs here: https://tools.ietf.org/html/rfc6570#section-2.

Expressions are placeholders which are to be substituted by the values their variables reference.

```text
http://example.org/~{username}/
http://example.org/dictionary/{term:1}/{term}
http://example.org/search{?q*,lang}
```

# How to use `uri`

This tool takes a JSON as input data model containing the values ​​of the variables. Example:

```json
{
    "username": "scarlett",
    "term": "black widow",
    "q": {
        "a": "mars",
        "b": "jupiter"
    },
    "lang": "en"
}
```

You can pass this JSON as file using the `-i` flag:

```bash
$ uri -i data.json http://example.org/~{username}/{term:1}/{term}{?q*,lang}
http://example.org/~scarlett/b/black%20widow?q*
```

Or you can pipe it directly to the tool (and this the most interesting use case):

```bash
$ cat data.json | uri http://example.org/~{username}/{term:1}/{term}{?q*,lang}
http://example.org/~scarlett/b/black%20widow?q*
```

Using commandline pipes you can achieve tasks like this:

```bash
$ cat testdata/pets.json \
  | jq --raw-output '.pets[1] | {category: .species, year: .birthYear}' \
  | uri https://pets.api.com/{category}{?year} \
  | xargs curl
```

# Installation Steps

In order to use the crumbs command, compile it using the following command:

```bash
go get -u github.com/lucasepe/uri
```

This will create the executable under your $GOPATH/bin directory.


## Ready-To-Use Releases 

If you don't want to compile the source code yourself, [here you can find the executables](https://github.com/lucasepe/uri/releases/latest) for:

- MacOS
- Linux
- Windows

## Credits

`uri` was possible thanks to the [Joshua Tacoma URI Template library](https://github.com/jtacoma/uritemplates).