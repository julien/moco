moco
----

*mocos* means "snot" in spanish, and this is
what you can probably expect from this code.

**moco** let's you "mock" http requests easily, you just need to pass a json file to the progam.

`moco -f FILENAME [-p PORT]`


The `-f` flag is mandatory and the `-p` flag is optional.

The json file should respect this structure:

*example*

```
{
  "/api/foo": {
    "headers": {
      "Content-type": "application/json",
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Headers": "X-Requested-With"
    },
    "statusCode": 200,
    "body": {
      title: "Some cool stuff",
      items: ["apple", "banana", "cherry", "watermelon"]
    }
  },

  "/api/bar": {
    "statusCode": 404,
    "body": "Get out of here"
  },

  "/api/baz/\\d{1,2}/profile": {
    "headers": {
      "Content-type": "application/json",
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Headers": "X-Requested-With",
      "X-App-Token": "b4c0n"
    },
    "body": {
      "message": "Welcome to my super API."
    }
  }
}
```

The only required field is the `body` but you can add `headers` and use a custom `statusCode` if you want to.

Regular expressions are supported although quite limited.


**NOTE**

This is a work in progress, which means you shouldn't expect
too much for now, that said you're more than welcome to contribute.

---

*If you're using [mocko](https://github.com/julien/mocko),
 note that this version uses less than 1/10th of RAM and does the same thing.*


LICENSE
-------

MIT
