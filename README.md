# doctrans-framework

The Document Transformation Application (DTA) is a microservice based web service application framework provided by the research group QDS of the TU Berlin.

It serves the research group [Quality of Distributed Systems](https://www.qds.tu-berlin.de/menue/quality_engineering_of_open_distributed_systems/) and [Open Distributed System](https://www.ods.tu-berlin.de/menue/fachgebiet_open_distributed_systems/) as example playground, to

- develop
- deploy, and
- test

micro-services and functions a service (FaaS) developments and deployment scenarios.

The DTA serves as a playground for different research topics and provides a collection of different document transformation functions.

A DocumentTransformationFunction DTF transforms a document into another document. Transforming in our context is meant in a very broad sense. Actually the semantics defined is, that one document, represented as a byte array, is passed as an argument to the DTF and another array of bytes is returned.
