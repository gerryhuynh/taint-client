## What

Emulates `kubectl taint nodes <nodename> dedicated=groupName:<NoSchedule | TaintEffect><remove>`

## Usage

Option 1: `go run main.go` to list nodes and associated taints

Option 2:

```
go run main.go
  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default $HOME)
  -node string
        Node name
  -remove
        Boolean to remove taint from node
  -taint string
        Taint Effect Options: NoSchedule, PreferNoSchedule, NoExecute (default "NoSchedule")
```

E.g.
`go run main.go -node node-123 -taint PreferNoSchedule -remove`

- Removes the taint `dedicated=groupNmae:PreferNoSchedule` from node `node-123`
