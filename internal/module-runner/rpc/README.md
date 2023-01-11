

### How to run the RPC

Have a look in: `internal/module-runner/rpc`

```go
// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"spark": &sparkv1.SparkPlugin{},
}

m, _ := yaml.Marshal(map[string]any{
    "config": map[string]string{
        "id":   s.Id,
        "name": s.Name,
    },
    "temporal": cfg.Temporal,
})
cmd := exec.Command(path.Join(cfg.BinBasePath, s.Name))
cmd.Env = os.Environ()
cmd.Env = append(cmd.Env, "SPARK_SECRET="+base64.StdEncoding.EncodeToString(m))

if s.Config != "" {
    cmd.Env = append(cmd.Env, "CONFIG_SECRET="+base64.StdEncoding.EncodeToString([]byte(s.Config)))
}

// We're a host! Start by launching the plugin process.
pc := plugin.NewClient(&plugin.ClientConfig{
    HandshakeConfig: plugin.HandshakeConfig{
        ProtocolVersion:  1,
        MagicCookieKey:   "BASIC_PLUGIN",
        MagicCookieValue: s.Id,
    },
    Plugins: pluginMap,
    Cmd:     cmd,
    Logger:  logger,
})

// Connect via RPC
rpcClient, err := pc.Client()
if err != nil {
    log.Fatal(err)
}

// Request the plugin
raw, err := rpcClient.Dispense("spark")
if err != nil {
    log.Fatal(err)
}

sa := raw.(sparkv1.SparkRpcApi)
log.Printf("pong: %s", sa.Greet(&sparkv1.IBlackboard{
    Value: "module runner",
    GetVal: func() string {
        return "its from here"
    },
}))
```

