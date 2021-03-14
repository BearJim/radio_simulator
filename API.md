# API

## CLI Commands

### get ues
- describe: get all the UEs in brief
- output example
```shell
SUPI                CM-STATE   RM-STATE      SERVING-RAN
imsi-2089300000003  IDLE       Registered    ran-1 (127.0.0.1:38412)
imsi-2089300000004  Connected  Registered    ran-1 (127.0.0.1:38412)
imsi-2089300000005  IDLE       Deregistered  ran-2 (127.0.0.2:38412)
```

### describe ue \<supi\>
- describe: describe UE context with *supi*
- output example
```shell
Supi: imsi-2089300000005
Guti: imsi-20893cafe0001
CM-State: IDLE
RM-State: Registered
Serving-RAN: ran-1 (127.0.0.1:38412)
RAN-UE-NGAP-ID: 1
AMF-UE-NGAP-ID: 1
```

### reg \<supi\> \<ranID\>
- describe: trigger initial registration procedure for UE with *supi* via RAN *ranID*

### dereg \<supi\>
- describe: trigger deregistration procedure for UE with *supi*

### idle \<supi\>
- describe: UE with *supi* will enter CM-IDLE state
    - AN Release procedure

### connect \<supi\>
- describe: UE with *supi* will enter CM-Connect state
    - UE-triggered Service request

---

### sess \<supi\> [add|delete] \<pduSessionID\>
- describe: add/delete pdu session for UE with *supi*
    - add: pdu session establishment
    - delete: pdu session release
