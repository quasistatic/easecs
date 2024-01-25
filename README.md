# EaseCS
CLI tool to quickly `exec` into any ECS service container of choice.

# Install

```shell
$ go install github.com/quasistatic/easecs
$ easecs help
```

# Usage

```shell
$ easecs exec <cluster> <service> [<container>]
```
The `cluster`, `service` and `container` args could be regexp strings like `prod$` which you 
know would match with the actual cluster name uniquely. Arbitrary option may be chosen in case 
of multiple matches. 

To quickly list what all you have in your ECS:
```shell
$ easecs list
```
