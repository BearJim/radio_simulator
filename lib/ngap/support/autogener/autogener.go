package ngap_autogener

import (
	"fmt"
	"os"
	"radio_simulator/lib/ngap/ngapType"
	"reflect"
	"strings"
)

func GenerateBuildTemp(present int, typePresent int, ies []int) string {

	out := `	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresent`

	critiMap := criticalityMap[present]
	pdu := ngapType.NGAPPDU{}
	presentName := presentMap[present]
	presentVariableName := strings.ToLower(presentName[:1]) + presentName[1:]
	value := reflect.TypeOf(pdu).Field(present).Type.Elem().Field(2).Type
	field := value.Field(typePresent)
	fieldName := field.Name
	fieldVariableName := strings.ToLower(fieldName[:1]) + fieldName[1:]
	fieldType := field.Type.String()
	fieldType = strings.Split(fieldType, ".")[1]

	out += fmt.Sprint(presentName, `
	pdu.`, presentName, ` = new(ngapType.`, presentName, `)

	`, presentVariableName, ` := pdu.`, presentName, `
	`, presentVariableName, `.ProcedureCode.Value = ngapType.ProcedureCode
	`, presentVariableName, `.Criticality.Value = ngapType.CriticalityPresent`, criticalityToString[critiMap[typePresent][0]], `
	
	`, presentVariableName, `.Value.Present = ngapType.`, presentName, `Present`, fieldName, `
	`, presentVariableName, `.Value.`, fieldName, ` = new(ngapType.`, fieldType, `)
	
	`, fieldVariableName, ` := `, presentVariableName, `.Value.`, fieldName, `
	`, fieldVariableName, `IEs := &`, fieldVariableName, `.ProtocolIEs`)

	fieldTypes := field.Type.Elem().Field(0).Type.Field(0).Type.Elem().Field(2).Type

	if ies == nil {
		for i := 1; i < fieldTypes.NumField(); i++ {
			ies = append(ies, i)
		}
	}

	for _, iePresent := range ies {
		ieType := fieldTypes.Field(iePresent)
		ieName := ieType.Name
		variableName := strings.ToLower(ieName[:1]) + ieName[1:]
		ieFieldType := ieType.Type.String()
		ieFieldType = strings.Split(ieFieldType, ".")[1]
		out += fmt.Sprint(`
	// `, ieName, `
	{
		ie := ngapType.`, fieldType, `IEs{}
		ie.Id.Value = ngapType.ProtocolIEID`, ieName, `
		ie.Criticality.Value = ngapType.CriticalityPresent`, criticalityToString[critiMap[typePresent][iePresent]], `
		ie.Value.Present = ngapType.`, fieldType, `IEsPresent`, ieName, `
		ie.Value.`, ieName, ` = new(ngapType.`, ieFieldType, `)

		`, variableName, ` := ie.Value.`, ieName, `
		
		`, fieldVariableName, `IEs.List = append(`, fieldVariableName, `IEs.List, ie)
	}`)
	}

	out += fmt.Sprintln("\n\n\treturn ngap.Encoder(pdu)")
	return out

}

func GenerateHandlerTemp(present int, typePresent int, ies []int) string {

	out := `

	if ran == nil{
		ngapLog.Error("RAN Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}
	`

	critiMap := criticalityMap[present]
	pdu := ngapType.NGAPPDU{}
	presentName := presentMap[present]
	presentVariableName := strings.ToLower(presentName[:1]) + presentName[1:]
	value := reflect.TypeOf(pdu).Field(present).Type.Elem().Field(2).Type
	field := value.Field(typePresent)
	fieldName := field.Name
	fieldVariableName := strings.ToLower(fieldName[:1]) + fieldName[1:]

	out += fmt.Sprint(`
	`, presentVariableName, ` := message.`, presentName, `
	if `, presentVariableName, ` == nil {
		ngapLog.Error("`, presentName, ` is nil")
		return
	}
	
	`, fieldVariableName, `:= `, presentVariableName, `.Value.`, fieldName, `
	if `, fieldVariableName, ` == nil {
		ngapLog.Error("`, fieldVariableName, ` is nil")
		return
	}

	for _, ie := range `, fieldVariableName, `.ProtocolIEs.List {
		switch ie.Id.Value {`)

	fieldTypes := field.Type.Elem().Field(0).Type.Field(0).Type.Elem().Field(2).Type

	if ies == nil {
		for i := 1; i < fieldTypes.NumField(); i++ {
			ies = append(ies, i)
		}
	}
	diagnostics := false
	prefix := ""
	for _, iePresent := range ies {
		ieType := fieldTypes.Field(iePresent)
		ieName := ieType.Name
		variableName := strings.ToLower(ieName[:1]) + ieName[1:]
		ieFieldType := ieType.Type.String()
		prefix += fmt.Sprintf("\tvar %s %s\n", variableName, ieFieldType)
		out += fmt.Sprint(`
		case ngapType.ProtocolIEID`, ieName, `:
			ngapLog.Traceln("[NGAP] Decode IE `, ieName, `")
			`, variableName, ` = ie.Value.`, ieName)
		if critiMap[typePresent][iePresent] == 1 {
			diagnostics = true
			out += fmt.Sprint(`
			if `, variableName, ` == nil {
				ngapLog.Error("`, ieName, ` is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}`)
		}
	}
	out += `
		}
	}
	`
	if diagnostics {
		out = "\n\tvar iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList" + out

		out += fmt.Sprint(`
	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCode
		trigger := ngapType.TriggeringMessagePresent`, presentName, `
		criticality := ngapType.CriticalityPresent`, criticalityToString[critiMap[typePresent][0]], `
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		simulator_ngap.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}`)

	}

	return fmt.Sprint(prefix, out)

}

func GenerateBuildTempALL(present int) {
	pdu := ngapType.NGAPPDU{}
	value := reflect.TypeOf(pdu).Field(present).Type.Elem().Field(2).Type
	for i := 1; i < value.NumField(); i++ {
		fileName := "build" + value.Field(i).Name + ".txt"
		f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		output := GenerateBuildTemp(present, i, nil)
		_, err := f.WriteString(output)
		if err != nil {
			fmt.Printf("write file [%s] error[%s]", fileName, err.Error())
		}
		f.Close()
	}
}

func GenerateHandlerTempALL(present int) {
	pdu := ngapType.NGAPPDU{}
	value := reflect.TypeOf(pdu).Field(present).Type.Elem().Field(2).Type
	for i := 1; i < value.NumField(); i++ {
		fileName := "handle" + value.Field(i).Name + ".txt"
		f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		output := GenerateHandlerTemp(present, i, nil)
		_, err := f.WriteString(output)
		if err != nil {
			fmt.Printf("write file [%s] error[%s]", fileName, err.Error())
		}
		f.Close()
	}
}
