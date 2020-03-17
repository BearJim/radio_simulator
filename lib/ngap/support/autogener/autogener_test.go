package ngap_autogener_test

import (
	"github.com/stretchr/testify/assert"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/lib/ngap/support/autogener"
	"os"

	"testing"
)

func TestAutogenerBuildModificationIndication(t *testing.T) {
	f, _ := os.OpenFile("build.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateBuildTemp(ngapType.NGAPPDUPresentInitiatingMessage, ngapType.InitiatingMessagePresentPDUSessionResourceModifyIndication, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}
func TestAutogenerHandlePDUModifyConfirm(t *testing.T) {
	f, _ := os.OpenFile("handler.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateHandlerTemp(ngapType.NGAPPDUPresentSuccessfulOutcome, ngapType.SuccessfulOutcomePresentPDUSessionResourceModifyConfirm, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerBuildPDUSessionRelRes(t *testing.T) {
	f, _ := os.OpenFile("build.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateBuildTemp(ngapType.NGAPPDUPresentSuccessfulOutcome, ngapType.SuccessfulOutcomePresentPDUSessionResourceReleaseResponse, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerHandlePDUReleaseCmd(t *testing.T) {
	f, _ := os.OpenFile("handler.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateHandlerTemp(ngapType.NGAPPDUPresentInitiatingMessage, ngapType.InitiatingMessagePresentPDUSessionResourceReleaseCommand, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerBuildPDUSessionNoti(t *testing.T) {
	f, _ := os.OpenFile("build.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateBuildTemp(ngapType.NGAPPDUPresentInitiatingMessage, ngapType.InitiatingMessagePresentPDUSessionResourceNotify, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerBuildNGReset(t *testing.T) {
	f, _ := os.OpenFile("build.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateBuildTemp(ngapType.NGAPPDUPresentInitiatingMessage, ngapType.InitiatingMessagePresentNGReset, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerBuildNGResetAck(t *testing.T) {
	f, _ := os.OpenFile("build.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateBuildTemp(ngapType.NGAPPDUPresentSuccessfulOutcome, ngapType.SuccessfulOutcomePresentNGResetAcknowledge, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerHandleNGReset(t *testing.T) {
	f, _ := os.OpenFile("handler.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateHandlerTemp(ngapType.NGAPPDUPresentInitiatingMessage, ngapType.InitiatingMessagePresentNGReset, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerHandleNGResetAck(t *testing.T) {
	f, _ := os.OpenFile("handler.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateHandlerTemp(ngapType.NGAPPDUPresentSuccessfulOutcome, ngapType.SuccessfulOutcomePresentNGResetAcknowledge, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}

func TestAutogenerHandleAll(t *testing.T) {
	ngap_autogener.GenerateHandlerTempALL(ngapType.NGAPPDUPresentInitiatingMessage)
	ngap_autogener.GenerateHandlerTempALL(ngapType.NGAPPDUPresentSuccessfulOutcome)
	ngap_autogener.GenerateHandlerTempALL(ngapType.NGAPPDUPresentUnsuccessfulOutcome)
}

func TestAutogenerBuildAll(t *testing.T) {
	ngap_autogener.GenerateBuildTempALL(ngapType.NGAPPDUPresentInitiatingMessage)
	ngap_autogener.GenerateBuildTempALL(ngapType.NGAPPDUPresentSuccessfulOutcome)
	ngap_autogener.GenerateBuildTempALL(ngapType.NGAPPDUPresentUnsuccessfulOutcome)
}
func TestAutogenerBuildInitiatingUeMessage(t *testing.T) {
	f, _ := os.OpenFile("build.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	output := ngap_autogener.GenerateBuildTemp(ngapType.NGAPPDUPresentInitiatingMessage, ngapType.InitiatingMessagePresentInitialUEMessage, nil)
	_, err := f.WriteString(output)
	assert.True(t, err == nil)
}
