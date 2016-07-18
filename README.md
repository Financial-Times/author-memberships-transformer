# author-memberships-transformer

[![CircleCI](https://circleci.com/gh/Financial-Times/author-memberships-transformer.svg?style=svg)](https://circleci.com/gh/Financial-Times/author-memberships-transformer)

Retrieves author data curated by editorial people and transforms it to People Memberships according to UP JSON model.
The service exposes endpoints for getting all the curated authors' membership UUIDs and for getting memberships by uuid.
The service consumes data specified by two Google Google spreadsheet, one contains authors' data and the another one contains authors roles.
Spreadsheet data is consumed by the service through [Bertha API](https://github.com/ft-interactive/bertha/wiki/Tutorial), which transforms Google spreadsheets data to JSON.
Output examples for authors and roles JSON from Berta are provided below.

####Bertha Authors
```
[
  {
    "name": "Martin Wolf",
    "role": "Columnist",
    "jobtitle": "Chief Economics Commentator",
    "email": "email@ft.com",
    "imageurl": "http//image.site.com/Martin_Wolf.png",
    "biography": "<p>Martin Wolf is chief economics commentator at the Financial Times, London. He was awarded the CBE (Commander of the British Empire) in 2000 “for services to financial journalism”.</p>",
    "twitterhandle": "@martinwolf_",
    "uuid": "daf5fed2-013c-468d-85c4-aee779b8aa54",
    "tmeidentifier": "Q0ItMDAwMDkwMA==-QXV0aG9ycw==",
    "membershipuuid": "03ca31b5-534a-4d26-bd09-3532b04640d5"
  },
  {
    "name": "Lucy Kellaway",
    "role": "Columnist",
    "jobtitle": "Associate Editor and Work & Career Columnist",
    "email": "email@ft.com",
    "imageurl": "http//image.site.com/Lucy_Kellaway.png",
    "biography": "<p>Lucy Kellaway is an Associate Editor and management columnist of the FT. For the past 15 years her weekly Monday column has poked fun at management fads and jargon and celebrated the ups and downs of office life.</p>",
    "twitterhandle": null,
    "uuid": "daf5fed2-013c-468d-85c4-aee779b8aa55",
    "tmeidentifier": "Q0ItMDAwMDkyNg==-QXV0aG9ycw==",
    "membershipuuid": "7318d8eb-5d48-4008-8744-aa3f18957216"
  },
  ...
]  
```

####Bertha Roles
```
[
  {
    "uuid": "33ee38a4-c677-4952-a141-2ae14da3aedd",
    "preflabel": "Journalist"
  },
  {
    "uuid": "7ef75a6a-b6bf-4eb7-a1da-03e0acabef1b",
    "preflabel": "Columnist"
  }
]
```

# How to run

## Locally: 

`go get github.com/Financial-Times/author-memberships-transformer`

`$GOPATH/bin/ ./author-memberships-transformer --bertha-authors-source-url=<BERTHA_AUTHORS_SOURCE_URL> --bertha-roles-source-url=<BERTHA_ROLES_SOURCE_URL> --port=8080`                

```
export|set PORT=8080
export|set BERTHA_AUTHORS_SOURCE_URL="http://.../Authors"
export|set BERTHA_ROLES_SOURCE_URL="http://.../Roles"
$GOPATH/bin/author-memberships-transformer
```

## With Docker:

`docker build -t coco/author-memberships-transformer .`

`docker run -ti --env BERTHA_AUTHORS_SOURCE_URL=<bertha_authors_url> --env BERTHA_ROLES_SOURCE_URL=<bertha_roles_url> coco/author-memberships-transformer`

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
This output data will be consumed as a stream by the [concept publisher](https://github.com/Financial-Times/concept-publisher).
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