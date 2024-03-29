**moco** let's you "mock" HTTP requests easily.

usage:

`moco [-f FILENAME] [-p PORT]`

Will start an HTTP server listening on the port specified by the `-p` flag or port 8000 by default.

Both the `-f` flag and the `-p` flag are optional.

You can use the `-p` flag to specify a JSON file that define "routes/endpoints" and it should have this structure:
```json
{
  "/api/foo": {
    "headers": {
      "Content-type": "application/json",
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Headers": "X-Requested-With"
    },
    "statusCode": 200,
    "body": {
      "title": "Some cool stuff",
      "items": ["apple", "banana", "cherry", "watermelon"]
    }
  },

  "/api/bar": {
    "statusCode": 404,
    "body": "Nothing found here"
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
All fields are optional.
Regular expressions are supported although quite limited.

If no flags are passed, `moco` will return a 200 status code for any request.

**NOTE**

This is a work in progress, which means you shouldn't expect
too much for now, that said you're more than welcome to contribute.

---

*If you're using [mocko](https://github.com/julien/mocko),
 note that this version uses less than 1/10th of RAM and does the same thing.*


LICENSE
-------

MIT
