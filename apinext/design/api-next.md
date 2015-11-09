
# Introduction

This API is intended to both harmonize what has previously been the fairly separate 'vision'
and 'curator' APIs, as well as lay a foundation for a fully enterprise-ready Clarifai API platform.

More information can be found in Google Drive in the [API Next folder](https://drive.google.com/a/clarifai.com/folderview?id=0B3i6Xc-tuTHlMzc1dldQYjA0VFE&usp=sharing)

# Key Features
A few of the new(er) features include

+ strongly typed models (input types and output prediction types)
+ extensible model types
+ discoverable models
+ ephemeral predictions not requiring any persistence in the platform
+ large batch (bulk) prediction requests
+ persisted predictions power search
+ customer-trained custom models act like regular system-provided models
+ No multiop (use model sets instead)
+ no facedetrec op (request a model that predicts bounding boxes for the concept 'face')

# Resource Model

The API resource model is the set of things the API lets the user manipulate. A strong goal here is
to use the same words that most users would use. We strive to avoid introducing objects or vocabulary
if users wouldn't naturally think of them. So, we use *image* and *video* to refer to those user data types.
The API offers services that make *predictions* using *models*, so those things are in the resource model.

+ Assets (aka user data types)
  + images
  + videos
  + audio
  + text
  + compound data types
  + Collections (of assets)
+ models
  + clarifai-provided 'standard' models i.e. the general model
  + clarifai-provided 'domain' models { nsfw, logos, e-commerce }
  + clarifai-trained 'customer' models i.e. models trained for specific customer use cases e.g. smartplate
  + user-trained custom models
+ prediction types
  + tags (will be deprecated soon, but supported until we fully retire those models)
  + concepts
  + embeddings
  + dominant colors
  + bounding boxes (by concept e.g. faces, logos)
+ Feedback
+ Users, Groups, Organization, plan
+ Features, Feature Flags

There are some more elaborated diagrams in the [Domain Model](https://drive.google.com/drive/folders/0B3i6Xc-tuTHlRmNFVUVmUmNfMTA) folder.

## About this Draft
Given our decision to document the API using Swagger, we don't need to cover every detail in a design document like this.
So, we intend to document examples of each type of endpoint, conventions, errors, headers only once
and then expect to continue those patterns into the Swagger spec. If we don't like what's here at the 'resource model'
or API spec 'sketch' level, there's little point in pursuing the next level of detail.

## Conventions
We try to follow generally accepted best practices. Of course, there's a fair amount of diversity of opinion on this. We
have to find our version of _best practices_. A couple of key references we lean on include

+ [best practices for a pragmatic restful api](http://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api)

Highlights:

+ JSON-only responses
+ Always pretty print responses, and enable GZIP
+ JSON encoded POST, PUT & PATCH bodies (require Content-Type header set to application/json or throw 415 Unsupported Media Type)
+ twitter-like, parse-like 'include=' for controlling returned content / depth

## Open Issues
+ resolve ids vs. names convention
+ hy-phens, under_scores, or camelCase [one discussion](http://stackoverflow.com/questions/10302179/hyphen-underscore-or-camelcase-as-word-delimiter-in-uris)
+ define the standard error payload, and the convention for endpoint-specific additional error payload
+ establish a convention for pagination of returned resources
+ establish convention for large bulk operations i.e. the bulk version of /predict
+ establish the convention for long-running operations
+ establish the convention for streaming operations

## Response Formats
We only provide responses in JSON format. We allow '.json' to be appended to any resource URI to ease compatibility with
clients that expect to provide the content type suffix.

We do not provide responses in XML, and appending '.xml' to resource URI will result in NOT FOUND errors.

## Exceptions
Like Stripe, we return conventional HTTP response codes to indicate the success or failure of an API request.
In general, codes in the ```2xx``` range indicate success, codes in the ```4xx``` range indicate an error that failed given the
information provided (e.g., a required parameter was omitted, a URL couldn't be accessed, etc.), and codes in the ```5xx``` range
indicate an error with Stripe's servers (these are rare).

Like Twilio, Clarifai returns exceptions in the HTTP response body when something goes wrong. An exception has the following
properties:

|Property|Description|
|--------|-----------|
|Status	|The HTTP status code for the exception.|
|Message|	A more descriptive message regarding the exception.|
|Service| The id/name of the service that encountered the error.|
|Code| An error code to find help for the exception.|
|MoreInfo| The URL of Clarifai's documentation for the error code.|
|Parameter (optional)|For request parameter errors, the name of the parameter|
|ParameterMessage (optional)|For requests parameter/payload errors, the specific problem encountered with the parameter|

If you receive an exception with status code 400 (invalid request), the 'Code' and 'MoreInfo' properties are useful for debugging what went wrong.

### Debug Mode (Clarifai internal)
In debug mode, the exception will include a source filename and line number and an optional stacktrace.

## Resource Types
The API exposes many resource types. Some correspond directly to 'objects' in the traditional sense and so are straightforward media representations of those objects. Some are *algorithmic resources* notably *models* and *query services*. Models make *predictions* of various types, and so they support a */predict* resource. The supported prediction types each have a media representation. *query* or *search* services also return their results in media representations suitable for their particular output. We try to pattern our algorithmic resources after common examples like [google search](www.google.com) and [Google Cloud Platform](https://cloud.google.com/prediction/docs/reference/v1.6/).

## User Data Types
The discussion of the resource model begins with the user's data types. In the current version of the API,
these are *image* and *video*. In the future, support will be extended to *audio*, *text*, and *compound*. Compound data types
are tuples of the primitive types e.g. ( audio, text ) or ( audio, video, text). Think of a video with an audio track and a subtitle track.

## Users
Note: this discussion of User is a placeholder while Vinay works on user-related API.
Every caller of the API is a User, and uses credentials to gain
access to the system. Users also own resources (images, models, etc.), and so many if not most objects have an owner
attribute that is a link to the User resource that owns it. User attributes include
+ id
+ username
+ full name
+ email address
+ dateJoined

TODO(jim,vinay) - vinay is working on a REST api for access the current users, applications, plans etc in the v1 platform.
We need to sort out a migration plan for supporting the richer v2 "user model".

## Models
Broadly, models represent entities that can examine media assets and make predictions
about the assets. Model attributes include

+ id
+ name
+ description
+ version
+ input-type
+ prediction-type
+ createdAt
+ owner

## Predictions
Broadly, *models* make *predictions* about *assets*. Models in the Clarifai platform each make a particular type of prediction
(sometimes more than one), and they know how to understand particular types of asset. So, we describe the *type* of a model by
specifying the type(s) of assets that it understands (its input-type) and the type of predictions that it makes (its
prediction-type).

Over time, the set of *model types* will grow, and the set of *models* will grow as the Clarifai platform gains new capabilities.

| Resource | Operation | Endpoint |Comments|
|----------|-----------|----------|--------|
|Model| List All Models | GET /models | system- and user-trained |
|Model| List Tag Models | GET /models?prediction-type=tag |  |
|Model| List Concept Models | GET /models?prediction-type=concept |  |
|Model| List Models that make predictions about images | GET /models?input-type=image |  |
|Model| List Models that predict concepts in images | GET /models?input-type=image,prediction-type=concept |  ||

### Input Types
Each model understands input of a particular type e.g. images or video or audio. The current input types are

+ image
+ video

todo: resolve how animated GIFs should be treated. While they share some traits of videos, there are semantic differences
that seem to argue for sugaring the APIs to treat animated GIFs as a sequence of images.

Input types to be supported in the future include:

+ text
+ audio
+ compound types

Where compound types are tuples of the base types e.g. (image,text) or (video, audio, text).


### Prediction Types
Currently, prediction type can have the following values.

+ tag
+ concept
+ embedding1024d
+ dominantColors
+ boundingBox

### List the Available Models
How can we get the list of models that we can use?

```bash
curl -X GET /models
```

The response body is a list of Model resources.
```js
"models" : [
    {
      "href" : "https://api.clarifai.com/v2/models/a7G4h6JU31"
      "createdAt" : "2011-11-07T20:58:34.448Z",
      "modelId": "a7G4h6JU31",
      "modelName": "model T",
      "modelDescription": "a basic model",
      "modelVersion" : "1.3",
      "inputType" : "image",
      "predictionType" : "tag",
      "ownerName" : "Clarifai Model Administrator",
      "owner" : "https://api.clarifai.com/v2/users/abf87F4EFs"
    },
    ...
  ]
}  
]
```

### Get the Latest Version of a given Model
Note that the models are versioned. There may be multiple versions of a single available at any time. If you want the predictions
returned by the platform to remain stable, you should request predictions from a particular modelId for the version you want to
use (and keep using).

If you want to use the latest version of a model without explicitly checking what the latest version is, there is a shorthand you
can use.

```bash
curl -X GET /models/:modelName/latest
```

The response body is the Model resource for the latest version of the model named **modelName**.
```js
{
  "href" : "https://api.clarifai.com/v2/models/a7G4h6JU31"
  "createdAt" : "2011-11-07T20:58:34.448Z",
  "modelId": "a7G4h6JU31",
  "modelName": "model T",
  "modelDescription": "a very snazzy model",
  "modelVersion" : "1.3",
  "inputType" : "image",
  "predictionType" : "tag",
  "ownerName" : "Clarifai Model Administrator",
  "owner" : "https://api.clarifai.com/v2/users/abf87F4EFs"
}
```

## Ephemeral Predictions
Our overall approach is dominated by asset management and curation use cases. There are some simple 'entry-level' use cases requiring
the platform to return predictions without storing references to media.
For now, we are calling these *ephemeral predictions* to express the the ephemeral nature of the interaction.

It's also convenient because this use case maps most directly to existing /v1/tag.

### Get the predicted Tags for an image from a specified Model

A Model is an _algorithmic resource_ (e.g. like a query service), and it offers a predict service to request
predictions for an image.

| Resource | Operation | Endpoint |Comments|
|----------|-----------|----------|--------|
|Model| Get Predictions | POST /models/:modelId/predict | ||

```bash
curl -X POST \
 -H "Authorization: Bearer <access_token>" \
 -H "Content-Type: application/json" \
 -H "Accept: application/json" \
 https://api.clarifai.com/v2/models/:modelId/predict
```

The request body is a list of uri that refer to publicly accessible images.

```js
{
  "uri" : [ "http://www.clarifai.com/img/metro-north.jpg", ... ]  
}
```

The response body is a list of the predictions made by the specified model for each of the specified images.
The type of resources returned is the model's prediction-type.

If the model's prediction type is **tag**, then the response body is a list of **tag-predictions**.

```js
{
  "status" : "ALL_OK"
  "status_message" : "All items in the request completed successfully"
  "predictions" : [
    {
      "uri" : "http://www.clarifai.com/img/metro-north.jpg",
      "modelId": "a7G4h6JU31",
      "tags" :
      [
        { "tag" : "dog", "confidence" : 0.88729 },
        { "tag" : "poodle", "confidence" : 0.74232 },
        { "tag" : "golden retriever", "confidence" : 0.32376 }
        ...
      ]          
    }
    ...
  ]
  "errors" : [
    {
      "uri" : "http://www.clarifai.com/img/metro-north.jpg",
      "model" : "https://api.clarifai.com/v2/models/uiSD78sm7V"
      "error" : {
          todo() specify the standard error body
      }
    }
    ...

  ]
}
```

The simplest possible way to get predictions is to GET the predictions for a
uri parameter.

| Resource | Operation | Endpoint |Comments|
|----------|-----------|----------|--------|
|Model| Get Predictions | GET /models/:modelId/predict?uri=:uriToPredict | ||

```bash
curl -X GET \
 -H "Authorization: Bearer <access_token>" \
 -H "Content-Type: application/json" \
 -H "Accept: application/json" \
 https://api.clarifai.com/v2/models/:modelId/predict?uri="http://www.clarifai.com/img/metro-north.jpg"
```

### Batch Request Status
When a caller requests prediction for multiple images in one request, the results may vary. The 'status' member
of the response body whether all, some, or none of the batch was successfully processed.

| Status | Description |
|--------|-------------|
|ALL_OK  | All of the items in the batch request completed successfully |
|ALL_ERROR | None of the items in the batch request completed successfully |
|PARTIAL_ERROR | Some of the items completed, some encountered errors |

### Tags
Tags are words that name objects in an image (e.g. dog, sailboat), or are descriptive of meaningful aspects of an image
(e.g. romantic, sunny).

We're moving quickly away from tags, but we keep them for now to exploit being able to call existing /v1 models that predict tags.

### What Tags will a (tag) Model Predict?
If a model predicts tags, we can ask the model to tell us what tags it will predict.

```bash
curl -X GET /models/:modelId/tags
```

The default response is the list of the tags that may be predicted by the model.

```js
"tags" : [
  "poodle",
  "golden retriever",
  "children with the head of a golden retriever"    
]
```

### What Concepts will a (concept) Model Predict?

```bash
curl -X GET /models/:modelId/concepts
```

The default response is the list of the concept (ids) that may be predicted by the model.

```js
"concepts": [
  {
    "conceptId" : "ai_HLmqFqBf"
    "localNames" : {
      "en" : "truck",
      "en-us" : "truck",
      "en-gb" : "lorry",
    }
  },
  ...
 ]
```

The concept ids can be used to query for the display names of the concepts.

```bash
curl -X GET /concepts/:conceptId
```

The default response is the list of the concepts that may be predicted by the model.

```js
"concepts": [
  {
    "conceptId" : "ai_HLmqFqBf"
    "localNames" : {
      "en" : "truck",
      "en-us" : "truck",
      "en-gb" : "lorry",
    }
  },
  ...
 ]
```

### Pagination

It's obvious that there will be too many resources to return in a single
response. We adopt the **offset,limit,total** convention, and we return
these attributes in the **meta** response field.

```js
"concepts": [
  {
    "conceptId" : "ai_HLmqFqBf"
    "localNames" : {
      "en" : "truck",
      "en-us" : "truck",
      "en-gb" : "lorry",
    }
  },
  ...
 ]
"meta" : {
  "paging" : {
    "offset" : 25,
    "limit" : 100,
    "total" : 11273
  }
}

```

### Getting Tags for Batches of Images
(refer to Parse scheme for bulk operations...)

Returns the tags predicted for the image. The returned tags include details about the model used to make the predictions. Note that since tag is a *prediction type* produced by arbitrarily many models, it's possible that the list of tags will be from more than one model. Similarly,

### No More 'default' Model
Models are versioned, and are specified using an id. So, it is unambiguous which version of a particular model you
are using.

As model versions are retired, requests that specify that model id will fail.
*Tag models* predict tags. As any given model evolves, the set of tags it predicts will change.
Consequently, the confidence of its predictions of any given class will change. So, model-v1.0
predict 'cat' with a confidence of 90%. model-v2.0 understands felines better, and so may start
to predict { 'cat', 'jaguar', 'leopard', 'lynx' }. With this change, the probability that 'cat' is
predicted for a particular image of a cat will change. If on Friday, the 'default' model is v1.0,
but on Monday it's moved to v2.0, your experience of asking for the predicted tags for your image of a
cat will change. Sometimes, this is what you want and sometimes it's not. To help you manage this, we
will let you ask for our 'latest' (and by implication, greatest) model of each type. But we'll also make it easy for you to ask for the same model day after day when that suits your needs better.

# Dominant Colors
A color extraction model will return information about dominant colors in an image or sequence of video frames.

# Custom Concept Models
tbd. for now, see a design discussion [here](https://docs.google.com/document/d/1tAlzIPdJ5fOa6_ZYgj-3SU1JcLYT5l3u-HxwuigKBH4/edit)

# Model Prediction Feedback

## User Metadata

# Advanced Topics

## Throttling

### Concurrent Request Limits
Our current throttling scheme doesn't explicitly limit a given client to a number of concurrent requests.
This might be a good idea.

+ Model Sets - Getting Predictions from multiple Models
+ Visual Search

## Getting Predictions from Multiple Models
Common use cases include: getting concept tags and the NSFW tag for an image, or getting the set
of predicted tags, and the bounding boxes if any for a specified set of tags.

This is a hybrid operation that doesn't map as nicely onto the pure resource model, but for
ease of use and efficiency reasons we don't want to force multiple calls per image.

```bash
curl -X GET /models/predict?models=model-id-1,model-id-2
```

# Parking Lot (Deferred Topics)
These things are in limbo - they've been brought up but really haven't evaluated or resolved in the context of API Next.

## Search based on Visual Similarity
One of the amazing features powered by the Clarifai API is visual search. Given a database of images D, provide a query image I_q, and return the set of images S that are the most visually similar to the
query image.

## Cluster Images by Visual Similarity
Closely related to searching by visual similarity, clustering will examine a database of images D and
will return k clusters of images where each cluster C_i will appear visually similar, and images in
each cluster will look more like each other than images in any other cluster.

## Model Sets
Although specifying a set of models with every operations is ok, we can make it more convenient and bring some other advantages if create a short-hand. For this, we'll create model set, which is just a list of models.

```bash
curl -X POST /modelsets/:modelsetId?models="general-v1.3,nsfw-latest"
curl -X POST /modelsets/:modelsetId?models="general-latest,logos-latest"
```

The advantage for you is you can set up the group once, and then continue to use 'standard' syntax.
The advantage for Clarifai is that by knowing which groups of models customers want to use (as
evidenced by these model set definitions) we can better optimize our compute resources.

## Predicting using Model Sets
---------------------------
One of the convenient things about model sets is they make it easy to get different types of predictions from different types of models in a single request. By now, you can probably predict what the endpoint will be for this operation. We don't show the result yet, but it will include the predictions from each model in the modelset.

```bash
curl -X POST /models/:modelsetId/predictions
```

## List Model Sets

```bash
curl -X GET /modelsets
```

## Delete Model Set

```bash
curl -X DELETE /modelsets/:modelsetId
```

## Accessing URIs that aren't publicly accessible

# Deferred
These things have been discussed and explicitly deferred. Effectively they are out of scope for this version.
