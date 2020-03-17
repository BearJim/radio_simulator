Get Expected Data
===========================================
Paste the ../15.2.0/38413.asn in https://asn1.io/asn1playground/ schema,
choose Encode and paste the Input Data (with specific Format), then click Encode,
you will get the expected Output (find data.per) for Encoded test

Input Data Example
===========================================
value NGAP-PDU ::= initiatingMessage{
    procedureCode  21,
    criticality  reject,
    value  NGSetupRequest : {
        protocolIEs {
            {
                id  27,
                criticality reject,
                value GlobalRANNodeID : globalGNB-ID{
                    pLMNIdentity '208f93'H,
                    gNB-ID gNB-ID : '454647'H
                    
                }
            },
            {
                id  82,
                criticality ignore,
                value RANNodeName : "free5GC"
            },
            {
                id  102,
                criticality reject,
                value SupportedTAList : {
                    {
                        tAC '001122'H,
                        broadcastPLMNList  {
                            {
                                pLMNIdentity '208f93'H,
                                tAISliceSupportList  {
                                    {
                                        s-NSSAI  {
                                            sST '01'H
                                            sD '010203'H
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            },
            {
                id  21,
                criticality ignore,
                value PagingDRX : 2
            }
            
        }
    }
}

Expected Output hex Data Example(remove space and newline manual)
===========================================
00150035 00000400 1B000800 208F9310 45464700 52400903 00667265 65354743
00660010 00000011 2200208F 93000010 08010203 00154001 40