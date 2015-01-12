moco
----

moco(s) means "snot" in spanish, and this is
what you can probably expect from this code.

moco let's you "mock" http requests easily,
you just need to pass a json file to the progam.

`moco -f FILENAME [-p PORT]`

The `-f` flag is mandatory and the `-p` flag is optional.

The json file should respect this structure:

```
{
  "/api/1": {

    "headers": {
      "Content-type": "application/json",
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Headers": "X-Requested-With"
    },

    "statusCode": 200,

    "body": {
      "name": "Something",
      "data": [1, 2, 3, 4, 5]
    }
  },
  "/api/2": {

    "headers": {
      "Content-type": "application/json",
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Headers": "X-Requested-With"
    },

    "statusCode": 201,

    "body": {
      "name": "Created",
      "data": {"user": "anonymous"}
    }
  }
}
```

The only required field is the `body` but you can add headers and
use custom status codes if you want to.

*If you're using [mocko](https://github.com/julien/mocko),
 note that this version uses less than 1/10th of RAM and does the same thing.*

LICENSE
-------

MIT
