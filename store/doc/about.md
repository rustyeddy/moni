---
title: About Store
description: Store is a Simple Object Store library useful for REST clients that receive and _temporarily_ store results.
categories: [ website ]
source: http://github.com/rustyeddy/store
site: http://rustyeddy/store
---

Store is a simple Object Storage library.  Directories are used as
storage _containers_ and files can be used to store _objects_.  

By default, store serializes/deserializes Go objects to/from JSON,
with some meta-info, including object Type.

Pretty much anytime of type of file can be _stored_ by _Store_.

Each _container_ has an index of objects in the _container_.

