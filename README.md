# imagepull-controller-workshop

## Agenda

### Introduction

- Explain the problem to solve.
  - https://github.com/titansoft-pte-ltd/imagepullsecret-patcher
  - https://github.com/titansoft-pte-ltd/imagepullsecret-patcher/blob/master/main.go#L110
  - https://medium.com/titansoft-engineering/kubernetes-cluster-wide-access-to-private-container-registry-with-imagepullsecret-patcher-b8b8fb79f7e5
- Explain controller components ([1], [2]).
  - ListerWatcher
  - Handler
  - Queue (ID and cache) and workers.
  - Optimizations made on the Listers.

### Implementation

- What is [Kooper] and alternatives.
  - [client-go].
  - [Kubebuilder].
  - [Operator framework][operator-framework].
- Explain project structure.
- Implement Retriever.
- Implement Handler.
- Test with local cluster using [kind].

### Optional homework

- Metrics.
- Unit testing.
- Secret.
  - Caching.
  - Update on change (Secret controller).
- Service account.
  - Update on change instead of ensuring.

### Further information

- Explain garbage collection (not required implementation).
  - Owner reference.
  - Finalizers.
- Creating an Operator.
  - Create CRDs with clients.
- Production ready operator example ([Bilrost]).

[1]: https://github.com/spotahome/gontroller
[2]: https://product.spotahome.com/gontroller-a-go-library-to-create-reliable-feedback-loop-controllers-832d4a9522ea
[kooper]: https://github.com/spotahome/kooper
[client-go]: https://github.com/kubernetes/client-go
[kubebuilder]: https://github.com/kubernetes-sigs/kubebuilder
[operator-framework]: https://github.com/operator-framework
[kind]: https://github.com/kubernetes-sigs/kind
[bilrost]: https://github.com/slok/bilrost
