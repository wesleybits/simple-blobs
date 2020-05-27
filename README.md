# Simple Blobs

This is a toy project I threw together to better understand how Akka Actors
changed. Also as an excuse to learn more about Akka HTTP. Also to play around
with the Cake pattern a little. Feel free to clone it to mess around with other
repository schemes, or use such specs laid out to mess with composing it with
other variations on the form.

This is a _very_ simple CRUD server done in a _very_ obtuse way. This is an
exercise.

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
