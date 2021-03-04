# imagepull-controller-workshop

## Versions

- [workshop]: First working version without some code for the workshop.
- [v0.1.0]: First working versions.
- [v0.2.0]: Optimized version for secret retrieval.

## Agenda

### Introduction

#### Problem

- Explain the problem to solve.
  - https://github.com/titansoft-pte-ltd/imagepullsecret-patcher
  - https://github.com/titansoft-pte-ltd/imagepullsecret-patcher/blob/master/main.go#L110
  - https://medium.com/titansoft-engineering/kubernetes-cluster-wide-access-to-private-container-registry-with-imagepullsecret-patcher-b8b8fb79f7e5

#### Concepts

- Controller Concepts
  - Edge and level trigger ([1][a], [2][b])
  - [Reconciliation]
- Explain controller components ([1][c], [2][d]).
  - ListerWatcher
  - Handler
  - Queue (ID and cache) and workers.
  - [Optimizations made on the Listers][loops-or-events].

### Implementation

#### Explanation

- What is [Kooper] and alternatives.
  - [client-go].
  - [Kubebuilder].
  - [Operator framework][operator-framework].
- Explain project structure.

#### Implement

```bash

git clone git@github.com:slok/imagepull-controller-workshop.git && cd ./imagepull-controller-workshop
git checkout workshop
```

- [Implement Retriever][workshop-retriever].
- [Implement Handler][workshop-handler].
- Test with local cluster using [kind].

```bash
kind create cluster
kubectl apply -f ./tests/manual/secret.yaml
go run ./cmd/imagepull-controller-workshop/ -r default --secret-name test-imagepull-credentials --development
```

### Optional homework

- Metrics.
- Unit testing.
- Secret.
  - Caching (Check [v0.2.0]).
  - Update on change (Secret controller).
- Service account.
  - Update on change instead of ensuring.

### Further information

- [Explain garbage collection][k8s-gc].
  - Owner reference.
  - Finalizers.
- Creating an Operator.
  - [kube-code-generator].
  - Create CRDs with clients.
    - Generate CRD ([1][e], [2][f], [3][g])
    - Generate Client ([1][h], [2][i])
- Production ready operator example ([Bilrost]).

[a]: https://speakerdeck.com/thockin/edge-vs-level-triggered-logic
[b]: https://hackernoon.com/level-triggering-and-reconciliation-in-kubernetes-1f17fe30333d
[c]: https://github.com/spotahome/gontroller
[d]: https://product.spotahome.com/gontroller-a-go-library-to-create-reliable-feedback-loop-controllers-832d4a9522ea
[e]: https://github.com/slok/kube-code-generator/blob/master/example/Makefile#L32
[f]: https://github.com/slok/kube-code-generator/blob/master/example/apis/comic/v1/types.go
[g]: https://github.com/slok/kube-code-generator/blob/master/example/manifests/comic.kube-code-generator.slok.dev_heroes.yaml
[h]: https://github.com/slok/kube-code-generator/blob/master/example/Makefile#L13
[i]: https://github.com/slok/kube-code-generator/blob/master/example/client/clientset/versioned/clientset.go#L32
[reconciliation]: https://speakerdeck.com/thockin/kubernetes-what-is-reconciliation
[kooper]: https://github.com/spotahome/kooper
[client-go]: https://github.com/kubernetes/client-go
[kubebuilder]: https://github.com/kubernetes-sigs/kubebuilder
[operator-framework]: https://github.com/operator-framework
[kind]: https://github.com/kubernetes-sigs/kind
[bilrost]: https://github.com/slok/bilrost
[loops-or-events]: https://speakerdeck.com/thockin/kubernetes-controllers-are-they-loops-or-events
[workshop]: https://github.com/slok/imagepull-controller-workshop/tree/workshop
[v0.1.0]: https://github.com/slok/imagepull-controller-workshop/tree/v0.1.0
[v0.2.0]: https://github.com/slok/imagepull-controller-workshop/tree/v0.2.0
[workshop-retriever]: https://github.com/slok/imagepull-controller-workshop/blob/workshop/internal/controller/retrieve.go#L21
[workshop-handler]: https://github.com/slok/imagepull-controller-workshop/blob/workshop/internal/controller/handle.go#L95-L101
[k8s-gc]: https://kubernetes.io/docs/concepts/workloads/controllers/garbage-collection/
[kube-code-generator]: https://github.com/slok/kube-code-generator
