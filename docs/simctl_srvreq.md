## simctl srvreq

trigger service request procedure for UE with SUPI

```
simctl srvreq <SUPI> [flags]
```

### Examples

```
srvreq imsi-2089300000001
```

### Options

```
  -f, --fail int32   trigger AMF fail ue count
  -h, --help         help for srvreq
```

### Options inherited from parent commands

```
      --db string   Database URL for simulator (default "mongodb://127.0.0.1:27017")
```

### SEE ALSO

* [simctl](simctl.md)	 - simctl - cli for Radio Simulator

###### Auto generated by spf13/cobra on 27-Jun-2021
