# blog-microservices

[![test](https://github.com/stonecutter/blog-microservices/actions/workflows/test.yaml/badge.svg)](https://github.com/stonecutter/blog-microservices/actions/workflows/test.yaml)

blog microservices deployed in an Istio-enabled kubernetes cluster.

### Architecture

![architecture](./assets/architecture.png)

### Full list what has been used:

* [gRPC](https://github.com/grpc/grpc-go) as transport layer
* [GORM](https://github.com/jackc/pgx) as database ORM
* [DTM](https://github.com/dtm-labs/dtm) as distributed transaction manager
* [Jaeger](https://www.jaegertracing.io/) open source, end-to-end distributed [tracing](https://opentracing.io/)
* [Prometheus](https://prometheus.io/) monitoring and alerting
* [Grafana](https://grafana.com/) for to compose observability dashboards with everything from Prometheus
* [Kiali](https://kiali.io/) The Console for Istio Service Mesh
* [Kubernetes](https://kubernetes.io/) for the cluster
* [Istio](https://istio.io/) as microservice architecture