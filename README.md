# Simple Blobs

These are very simple REST-y CRUD services written in a variety of languages. I
use this as an exercise to better understand how to to do 2 things:

- HTTP service with minimal external dependencies (aside from Scala)
- Module decoupling -- to some extent

Currently Simple Globs is implemented in the following languages:
- Scala2 with Akka and Circe
- Golang with some UUID lib from Google

All of these are intended to present the same REST interface, which just exposes
a super-simplified in-memory key/value storage with no mind paid for memory
efficiency or persistance.

## What it does

This is a simple service that receives data blobs and stores them. They're
assigned an UUID, and whenever any of them change, their hooks are called.

The general schema used throughout the project is:
```json
{
  "data": "formless JSON blob; anything you want",
  "hooks": [
    "array of URIs to POST change information to."
  ],
  "id": "system managed, it's a UUID that's has no other significance"
}
```

This thing has dynamic webhooks. You should be prepared to receive the following
at your webhook endpoints:

```json
{
  "data": {
    "//": "same as above; so repeated bellow",
    "data": "json blob",
    "hooks": [ "webhooks" ],
    "id": "system-assigned UUID"
  },
  "change": " Create | Update | Delete "
}
```

## Running It

Get SBT and do:

```sh
$ sbt run
```

## Endpoints

```
POST   /items      -> create an item
GET    /items      -> get all items
GET    /item/:uuid -> get one item with the matching UUID
PUT    /item/:uuid -> change the item under that UUID
DELETE /item/:uuid -> delete it
```

`GET` and `POST` endpoints call webhooks. `PUT` calls all webhooks in both the
old and new data blob.

## Does it do queries?

No.
