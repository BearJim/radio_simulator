package simulator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/free5gc/MongoDBLibrary"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
	"github.com/jay16213/radio_simulator/pkg/api"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/ue_factory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type ranApiURL struct {
	Name string `bson:"name"`
	Url  string `bson:"url"`
}

type Simulator struct {
	cc       *exec.Cmd
	dbClient *MongoDBLibrary.Client
}

func New(dbName string, dbUrl string) (*Simulator, error) {
	client, err := MongoDBLibrary.New(dbName, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("init DB error: %+v", err)
	}

	s := &Simulator{
		// RanPool: make(map[string]api.APIServiceClient),
		// UeContextPool: make(map[string]*simulator_context.UeContext),
		dbClient: client,
	}

	// cur, err := s.dbClient.Database().Collection("ranAPIList").Find(context.TODO(), bson.D{})
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	return nil, fmt.Errorf("init DB error: %+v", err)
	// }
	// defer cur.Close(context.TODO())
	// for cur.Next(context.TODO()) {
	// 	var ran ranApiURL
	// 	err := cur.Decode(&ran)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("decode url error: %+v\n", err)
	// 	}
	// 	if _, err := s.ConnectToRAN(ran.url); err != nil {
	// 		fmt.Printf("connect %s error: %+v", ran.url, err)
	// 		_, err := s.dbClient.Database().Collection("ranAPIList").DeleteOne(context.TODO(), bson.M{"url": ran.url})
	// 		if err != nil {
	// 			fmt.Printf("DeleteOne error: %+v", err)
	// 		}
	// 	}
	// }
	return s, nil
}

func (s *Simulator) StartNewRan() {
	c := exec.Command("./bin/simulator")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	if err := c.Start(); err != nil {
		fmt.Printf("run error: %+v\n", err)
	} else {
		fmt.Printf("c.Run err is nil\n")
	}
	s.cc = c
}

func (s *Simulator) GetRANs() {
	// for ranName, ranClient := range s.RanPool {
	// 	resp, err := ranClient.DescribeRAN(context.TODO(), &api.DescribeRANRequest{})
	// 	if err != nil {
	// 		fmt.Printf("fetch %s error: %+v\n", ranName, err)
	// 		continue
	// 	}

	// 	fmt.Printf("resp: %+v", resp.Name)
	// }
}

// ParseUEData read UE contexts from files then return a slice of *UeContext
func (s *Simulator) ParseUEData(filePath string) []*simulator_context.UeContext {
	var ueContexts []*simulator_context.UeContext

	ue := ue_factory.InitUeContextFactory(filePath)
	switch ue.IntegrityAlgStr {
	case "NIA0":
		ue.IntegrityAlg = security.AlgIntegrity128NIA0
	case "NIA1":
		ue.IntegrityAlg = security.AlgIntegrity128NIA1
	case "NIA2":
		ue.IntegrityAlg = security.AlgIntegrity128NIA2
	case "NIA3":
		ue.IntegrityAlg = security.AlgIntegrity128NIA3
	}

	switch ue.CipheringAlgStr {
	case "NEA0":
		ue.CipheringAlg = security.AlgCiphering128NEA0
	case "NEA1":
		ue.CipheringAlg = security.AlgCiphering128NEA1
	case "NEA2":
		ue.CipheringAlg = security.AlgCiphering128NEA2
	case "NEA3":
		ue.CipheringAlg = security.AlgCiphering128NEA3
	}
	ueContexts = append(ueContexts, ue)
	return ueContexts
}

func (s *Simulator) InsertUEContextToDB(ueContexts []*simulator_context.UeContext) {
	for _, ue := range ueContexts {
		s.updateUE(ue)
	}
}

func (s *Simulator) GetUEs() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "SUPI\tCM-STATE\tRM-STATE\tSERVING-RAN")

	cur, err := s.dbClient.Database().Collection("ue").Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Printf("fetch ue error: %+v\n", err)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var ue simulator_context.UeContext
		err := cur.Decode(&ue)
		if err != nil {
			fmt.Printf("decode ue error: %+v\n", err)
			continue
		}
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", ue.Supi, ue.CmState, ue.RmState, ue.ServingRan)
	}
	writer.Flush()
}

func (s *Simulator) DescribeUE(supi string) {
	ue, err := s.findUE(supi)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := s.connectToRan(ue.ServingRan)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.DescribeUE(context.TODO(), &api.DescribeUERequest{Supi: ue.Supi})
	if err != nil {
		fmt.Println(err)
		return
	}

	ueCtx := resp.GetUeContext()
	// ue.Guti
	ue.AmfUeNgapId = ueCtx.AmfUeNgapId
	ue.RanUeNgapId = ueCtx.RanUeNgapId
	ue.DLCount = security.Count(ueCtx.NasDownlinkCount)
	ue.ULCount = security.Count(ueCtx.NasUplinkCount)
	ue.RmState = ueCtx.RmState
	ue.CmState = ueCtx.CmState

	s.updateUE(ue)

	fmt.Printf("SUPI: %s\n", ue.Supi)
	fmt.Printf("AmfUeNgapId: %d\n", ue.AmfUeNgapId)
	fmt.Printf("RanUeNgapId: %d\n", ue.RanUeNgapId)
	fmt.Printf("DLCount: %d\n", ue.DLCount)
	fmt.Printf("ULCount: %d\n", ue.ULCount)
	fmt.Printf("RmState: %s\n", ue.RmState)
	fmt.Printf("CmState: %s\n", ue.CmState)
}

func (s *Simulator) UeRegister(supi string, ranName string) {
	ue, err := s.findUE(supi)
	if err != nil {
		fmt.Println(err)
		return
	}

	apiClient, err := s.connectToRan(ranName)
	if err != nil {
		fmt.Println(err)
		return
	}

	ue.ServingRan = ranName

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	startTime := time.Now()
	regResult, err := apiClient.Register(ctx, &api.RegisterRequest{
		Supi:         ue.Supi,
		ServingPlmn:  ue.ServingPlmnId,
		CipheringAlg: ue.CipheringAlgStr,
		IntegrityAlg: ue.IntegrityAlgStr,
		AuthMethod:   ue.AuthData.AuthMethod,
		K:            ue.AuthData.K,
		Opc:          ue.AuthData.Opc,
		Op:           ue.AuthData.Op,
		Amf:          ue.AuthData.AMF,
		Sqn:          ue.AuthData.SQN,
	})
	if err != nil {
		fmt.Printf("register error: %+v\n", err)
		return
	}

	finishTime := time.Since(startTime)

	if regResult.GetStatusCode() == api.StatusCode_ERROR {
		fmt.Printf("registration start failed: %s\n", regResult.GetBody())
	} else {
		fmt.Printf("registration success\n")
		resultUe := regResult.GetUeContext()
		ue.RmState = resultUe.GetRmState()
		ue.CmState = resultUe.GetCmState()
		ue.AmfUeNgapId = resultUe.GetAmfUeNgapId()
		ue.RanUeNgapId = resultUe.GetRanUeNgapId()
		ue.DLCount = security.Count(resultUe.GetNasDownlinkCount())
		ue.ULCount = security.Count(resultUe.GetNasUplinkCount())
		fmt.Printf("expand %+v\n", finishTime)
	}

	// update SQN when trigger registration
	num, _ := strconv.ParseInt(ue.AuthData.SQN, 16, 64)
	ue.AuthData.SQN = fmt.Sprintf("%x", num+1)
	s.updateUE(ue)
}

func (s *Simulator) UeDeregister(supi string) {
	ue, err := s.findUE(supi)
	if err != nil {
		fmt.Println(err)
		return
	}

	apiClient, err := s.connectToRan(ue.ServingRan)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	result, err := apiClient.Deregister(ctx, &api.DeregisterRequest{Supi: ue.Supi})
	if err != nil {
		fmt.Printf("deregister error: %+v\n", err)
		return
	}
	if result.GetStatusCode() == api.StatusCode_ERROR {
		fmt.Printf("Deregistration start failed: %s\n", result.GetBody())
	} else {
		fmt.Println("Deregistration success")
	}
}

func (s *Simulator) SubscribeUELog(client api.APIServiceClient, ue *simulator_context.UeContext, closeMsg []string) {
	stream, err := client.SubscribeLog(context.TODO(), &api.LogStreamingRequest{Supi: ue.Supi})
	if err != nil {
		fmt.Printf("subscribe ue log error: %+v\n", err)
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("close ue log streaming")
					return
				} else {
					fmt.Printf("recv ue log error: %+v\n", err)
					return
				}
			}
			fmt.Println(resp.LogMessage)
			// for _, msg := range closeMsg {
			// 	if resp.LogMessage == msg {
			// 		if err := stream.CloseSend(); err != nil {
			// 			fmt.Printf("CloseSend error: %+v\n", err)
			// 		}
			// 		fmt.Println("close log streaming")
			// 		return
			// 	}
			// }
		}
	}()
}

func (s *Simulator) UploadUEProfile(dbName string, dbUrl string) {
	// connect to free5gc DB
	dbClient, err := MongoDBLibrary.New(dbName, dbUrl)
	if err != nil {
		fmt.Printf("connect db error: %+v\n", err)
		return
	}

	// find all UE and upload to free5gc database
	cur, err := s.dbClient.Database().Collection("ue").Find(context.TODO(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No UE found in DB")
		} else {
			fmt.Printf("Find UE error: %+v\n", err)
		}
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var ue simulator_context.UeContext
		err := cur.Decode(&ue)
		if err != nil {
			fmt.Printf("decode ue error: %+v\n", err)
			continue
		}

		amData := models.AccessAndMobilitySubscriptionData{
			Gpsis: ue.Gpsis,
			Nssai: &ue.Nssai,
		}
		amPolicy := models.AmPolicyData{
			SubscCats: ue.SubscCats,
		}
		auths := ue.AuthData
		authsSubs := models.AuthenticationSubscription{
			AuthenticationMethod:          models.AuthMethod(auths.AuthMethod),
			AuthenticationManagementField: auths.AMF,
			PermanentKey: &models.PermanentKey{
				PermanentKeyValue: auths.K,
			},
			SequenceNumber: auths.SQN,
			Milenage:       &models.Milenage{},
		}
		if auths.Opc != "" {
			authsSubs.Opc = &models.Opc{
				OpcValue: auths.Opc,
			}
		}
		if auths.Op != "" {
			authsSubs.Milenage.Op = &models.Op{
				OpValue: auths.Op,
			}
		}
		InsertAuthSubscriptionToMongoDB(dbClient, ue.Supi, authsSubs)
		InsertAccessAndMobilitySubscriptionDataToMongoDB(dbClient, ue.Supi, amData, ue.ServingPlmnId)
		InsertSmfSelectionSubscriptionDataToMongoDB(dbClient, ue.Supi, ue.SmfSelData, ue.ServingPlmnId)
		InsertAmPolicyDataToMongoDB(dbClient, ue.Supi, amPolicy)
	}
}

func (s *Simulator) connectToRan(ranName string) (api.APIServiceClient, error) {
	result := s.dbClient.Database().Collection("ran").FindOne(context.TODO(), bson.M{"name": ranName})
	if result == nil || result.Err() == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("RAN not found (name: %s)", ranName)
	}
	var ranApi ranApiURL
	if err := result.Decode(&ranApi); err != nil {
		return nil, fmt.Errorf("decode ue error: %+v\n", err)
	}

	fmt.Printf("Connect to %s (%s)...\n", ranApi.Name, ranApi.Url)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, ranApi.Url, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New("connection timeout")
		} else {
			return nil, err
		}
	}
	return api.NewAPIServiceClient(conn), nil
}

func (s *Simulator) findUE(supi string) (*simulator_context.UeContext, error) {
	result := s.dbClient.Database().Collection("ue").FindOne(context.TODO(), bson.M{"supi": supi})
	if result == nil || result.Err() == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("UE not found (supi: %s)", supi)
	}
	var ue simulator_context.UeContext
	if err := result.Decode(&ue); err != nil {
		return nil, fmt.Errorf("decode ue error: %+v\n", err)
	}
	return &ue, nil
}

func (s *Simulator) updateUE(ue *simulator_context.UeContext) {
	upsert := true
	_, err := s.dbClient.Database().Collection("ue").UpdateOne(context.Background(), bson.M{"supi": ue.Supi},
		bson.M{"$set": ue}, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		fmt.Printf("update UE error: %+v\n", err)
	} else {
		fmt.Printf("UE profile updated (supi: %s)\n", ue.Supi)
	}
}

func InsertAuthSubscriptionToMongoDB(client *MongoDBLibrary.Client, ueId string, authSubs models.AuthenticationSubscription) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}
	putData := toBsonM(authSubs)
	putData["ueId"] = ueId
	client.RestfulAPIPutOne(collName, filter, putData)
}

func InsertAccessAndMobilitySubscriptionDataToMongoDB(client *MongoDBLibrary.Client, ueId string, amData models.AccessAndMobilitySubscriptionData, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(amData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	client.RestfulAPIPutOne(collName, filter, putData)
}

func InsertSmfSelectionSubscriptionDataToMongoDB(client *MongoDBLibrary.Client, ueId string, smfSelData models.SmfSelectionSubscriptionData, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(smfSelData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	client.RestfulAPIPutOne(collName, filter, putData)
}

func InsertAmPolicyDataToMongoDB(client *MongoDBLibrary.Client, ueId string, amPolicyData models.AmPolicyData) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}
	putData := toBsonM(amPolicyData)
	putData["ueId"] = ueId
	client.RestfulAPIPutOne(collName, filter, putData)
}

func toBsonM(data interface{}) bson.M {
	tmp, _ := json.Marshal(data)
	var putData = bson.M{}
	_ = json.Unmarshal(tmp, &putData)
	return putData
}