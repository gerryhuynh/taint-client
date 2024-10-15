## What

Emulates `kubectl taint nodes <nodename> dedicated=groupName:<NoSchedule | TaintEffect><remove>`

## Usage

```
go run main.go
  -node string
        (required) Node name
  -remove
        Boolean to remove taint from node
  -taint string
        (default "NoSchedule") Taint Effect Options: NoSchedule, PreferNoSchedule, NoExecute
```

E.g.
`go run main.go -node node-123 -taint PreferNoSchedule -remove`

- Removes the taint `dedicated=groupName:PreferNoSchedule` from node `node-123`
