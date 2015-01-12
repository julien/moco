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

  "/api/baz/qux": {
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

### TODO

- [ ] Watch .json file for changes and reload routes.
- [ ] Support regular experssion for routes.


**NOTE**

This is a work in progress, which means you shouldn't expect
too much for now, that said you're more than welcome to contribute.

---

*If you're using [mocko](https://github.com/julien/mocko),
 note that this version uses less than 1/10th of RAM and does the same thing.*


LICENSE
-------

This software is licensed under the MIT License.

Copyright Julien Castelain, 2015.

Permission is hereby granted, free of charge, to any person obtaining a copy of this
software and associated documentation files (the "Software"), to deal in the Software
without restriction, including without limitation the rights to use, copy, modify, merge,
publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or
substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.

