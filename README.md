# curated-authors-membership-transformer

[![CircleCI](https://circleci.com/gh/Financial-Times/curated-authors-membership-transformer.svg?style=svg)](https://circleci.com/gh/Financial-Times/curated-authors-membership-transformer)

Retrieves author data curated by editorial people and transforms it to People Membership according to UP JSON model.
Authors data is specified by a Google spreadsheet which is accessible through [Bertha API](https://github.com/ft-interactive/bertha/wiki/Tutorial).
The service exposes endpoints for getting all the curated authors' UUIDs and for getting authors by uuid.

# How to run

## Locally: 

`go get github.com/Financial-Times/curated-authors-membership-transformer`

`$GOPATH/bin/ ./curated-authors-membership-transformer --bertha-authors-source-url=<BERTHA_AUTHORS_SOURCE_URL> --bertha-roles-source-url=<BERTHA_ROLES_SOURCE_URL> --port=8080`                

```
export|set PORT=8080
export|set BERTHA_AUTHORS_SOURCE_URL="http://.../Authors"
export|set BERTHA_ROLES_SOURCE_URL="http://.../Roles"
$GOPATH/bin/curated-authors-membership-transformer
```

## With Docker:

`docker build -t coco/curated-authors-membership-transformer .`

`docker run -ti --env BERTHA_AUTHORS_SOURCE_URL=<bertha_authors_url> --env BERTHA_ROLES_SOURCE_URL=<bertha_roles_url> coco/curated-authors-membership-transformer`

#Endpoints

##Count
`GET /transformers/author-memberships/__count` returns the number of available memberships to be transformed as plain text.
A response example is provided below.

```
2
```

##IDs
`GET /transformers/author-memberships/__ids` returns the list of membership's UUIDs available to be transformed. 
The output is a sequence of JSON objects, however the returned `Content-Type` header is `text\plain`.
A response example is provided below.

```
{"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd2"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd5"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd9"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd8"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd0"} {"id":"daf5fed2-013c-468d-85c4-aee779b8aa53"} {"id":"daf5fed2-013c-468d-85c4-aee779b8aa51"} 
```

##Authors by UUID
`GET /transformers/author-memberships/{uuid}` returns author membership data of the given uuid. 
A response example is provided below.

```
{
  "uuid": "e6e8b0e4-4833-11e6-beb8-9e71128cae77",
  "prefLabel": "Chief Economics Commentator",
  "personUuid": "6f53299a-799d-49ae-ae9d-7e1f298daef7",
  "organisationUuid": "dac01f07-4b6d-3615-8532-a56752cc7e5f",
  "alternativeIdentifiers": {
    "uuids": [
      "e6e8b0e4-4833-11e6-beb8-9e71128cae77"
    ]
  },
  "membershipRoles": [
    {
      "roleuuid": "7ef75a6a-b6bf-4eb7-a1da-03e0acabef1b"
    }
  ]
}
```