package simulator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/free5gc/MongoDBLibrary"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
	"github.com/jay16213/radio_simulator/pkg/api"
	"github.com/jay16213/radio_simulator/pkg/logger"
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

func (s *Simulator) AllUeRegister(ranName string, triggerFailCount int, followOnRequest bool) {
	apiClient, err := s.connectToRan(ranName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// find all UE and register
	cur, err := s.dbClient.Database().Collection("ue").Find(context.TODO(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No UE found in DB")
		} else {
			fmt.Printf("Find UE error: %+v\n", err)
		}
		return
	}

	var ues []*simulator_context.UeContext
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		ue := new(simulator_context.UeContext)
		err := cur.Decode(ue)
		if err != nil {
			fmt.Printf("decode ue error: %+v\n", err)
			continue
		}
		ue.ServingRan = ranName
		ue.FollowOnRequest = followOnRequest
		ues = append(ues, ue)
	}

	// set fail triggering
	if triggerFailCount != 0 {
		triggerAmfFail(triggerFailCount)
	}

	successCnt := uint32(0)
	restartCnt := uint32(0)

	uePerSecond := 30
	if val, ok := os.LookupEnv("THESIS_UE_PER_SECOND"); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Parse THESIS_UE_PER_SECOND error: %+v", err)
			return
		}
		uePerSecond = v
	}
	timeSlot := 100 * time.Millisecond
	if val, ok := os.LookupEnv("THESIS_REQUEST_TIME_SLOT"); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Parse THESIS_REQUEST_TIME_SLOT error: %+v", err)
			return
		}
		timeSlot = time.Duration(v) * time.Millisecond
	}

	registrationAttempt := 3
	if val, ok := os.LookupEnv("THESIS_REGIS_ATTEMPT"); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Parse THESIS_REGIS_ATTEMPT error: %+v", err)
			return
		}
		registrationAttempt = v
	}

	wg := sync.WaitGroup{}
	startTime := time.Now()
	// fmt.Printf("%+v\n", startTime)
	for i := range ues {
		wg.Add(1)
		if i != 0 && i%uePerSecond == 0 {
			time.Sleep(timeSlot)
		}
		go func(ue *simulator_context.UeContext, wg *sync.WaitGroup) {
			success, now, completeTime, redoTime := s.ueRegister(ue, apiClient, registrationAttempt)
			if success {
				if redoTime != nil {
					fmt.Printf("%s, %+v, %d, %d\n", ue.Supi, now.Sub(startTime).Milliseconds(), completeTime.Milliseconds(),
						redoTime.Milliseconds())
					atomic.AddUint32(&restartCnt, 1)
				} else {
					// supi, time from system first reg, time from first reg this UE sent
					fmt.Printf("%s, %+v, %d\n", ue.Supi, now.Sub(startTime).Milliseconds(), completeTime.Milliseconds())
				}
				atomic.AddUint32(&successCnt, 1)
			}
			wg.Done()
		}(ues[i], &wg)
	}
	wg.Wait()
	fmt.Printf("%d, %d\n", successCnt, restartCnt)

	for _, ue := range ues {
		// update SQN when triggering registration
		ue.AuthDataSQNAddOne()
		s.updateUE(ue)
	}
}

func (s *Simulator) SingleUeRegister(supi string, ranName string, triggerFailCount int, followOnRequest bool) {
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

	// trigger fail
	if triggerFailCount != 0 {
		triggerAmfFail(triggerFailCount)
	}

	ue.ServingRan = ranName
	ue.FollowOnRequest = followOnRequest
	successCnt := uint32(0)
	restartCnt := uint32(0)
	wg := sync.WaitGroup{}
	wg.Add(1)

	startTime := time.Now()
	fmt.Printf("%+v\n", startTime)
	go func(ue *simulator_context.UeContext, wg *sync.WaitGroup) {
		success, now, completeTime, redoTime := s.ueRegister(ue, apiClient, 0)
		if success {
			if redoTime != nil {
				fmt.Printf("%s, %+v, %d, %d\n", ue.Supi, now, completeTime.Milliseconds(),
					redoTime.Milliseconds())
				atomic.AddUint32(&restartCnt, 1)
			} else {
				fmt.Printf("%s, %+v, %d\n", ue.Supi, now, completeTime.Milliseconds())
			}
			atomic.AddUint32(&successCnt, 1)
		}
		wg.Done()
	}(ue, &wg)
	wg.Wait()
}

func (s *Simulator) SingleUeServiceRequest(supi string, triggerFailCount int) {
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

	// trigger fail
	if triggerFailCount != 0 {
		triggerAmfFail(triggerFailCount)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		s.ueServiceRequest(ue, apiClient)
		wg.Done()
	}(&wg)
	wg.Wait()
}

func (s *Simulator) AllUeServiceRequest(ranName string, supis []string, triggerFailCount int) {
	apiClient, err := s.connectToRan(ranName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// trigger fail
	if triggerFailCount != 0 {
		triggerAmfFail(triggerFailCount)
	}

	uePerSecond := 50
	if val, ok := os.LookupEnv("THESIS_UE_PER_SECOND_SRV"); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Parse THESIS_UE_PER_SECOND_SRV error: %+v", err)
			return
		}
		uePerSecond = v
	}
	timeSlot := 100 * time.Millisecond
	if val, ok := os.LookupEnv("THESIS_REQUEST_TIME_SLOT_SRV"); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Parse THESIS_REQUEST_TIME_SLOT_SRV error: %+v", err)
			return
		}
		timeSlot = time.Duration(v) * time.Millisecond
	}
	successCnt := uint32(0)
	restartCnt := uint32(0)

	// find all UE and register
	cur, err := s.dbClient.Database().Collection("ue").Find(context.TODO(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No UE found in DB")
		} else {
			fmt.Printf("Find UE error: %+v\n", err)
		}
		return
	}

	var ues []*simulator_context.UeContext
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		ue := new(simulator_context.UeContext)
		err := cur.Decode(ue)
		if err != nil {
			fmt.Printf("decode ue error: %+v\n", err)
			continue
		}
		for _, supi := range supis {
			if ue.Supi == supi {
				ues = append(ues, ue)
				break
			}
		}
	}

	wg := sync.WaitGroup{}
	startTime := time.Now()
	for i := range ues {
		wg.Add(1)
		if i != 0 && i%uePerSecond == 0 {
			time.Sleep(timeSlot)
		}
		go func(ue *simulator_context.UeContext, wg *sync.WaitGroup) {
			success, now, completeTime, redoTime := s.ueServiceRequest(ue, apiClient)
			if success {
				if redoTime != nil {
					fmt.Printf("%s, %+v, %d, %d\n", ue.Supi, now.Sub(startTime).Milliseconds(), completeTime.Milliseconds(),
						redoTime.Milliseconds())
					atomic.AddUint32(&restartCnt, 1)
				} else {
					fmt.Printf("%s, %+v, %d\n", ue.Supi, now.Sub(startTime).Milliseconds(), completeTime.Milliseconds())
				}
				atomic.AddUint32(&successCnt, 1)
			}
			wg.Done()
		}(ues[i], &wg)
	}
	wg.Wait()
	fmt.Printf("%d, %d\n", successCnt, restartCnt)

	for _, ue := range ues {
		s.updateUE(ue)
	}
}

func (s *Simulator) ueServiceRequest(ue *simulator_context.UeContext, apiClient api.APIServiceClient) (
	success bool, nowTime time.Time, completeTime time.Duration, redoTime *time.Duration) {
	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()

	startTime := time.Now()
	result, err := apiClient.ServiceRequestProc(ctx, &api.ServiceRequest{
		Supi:        ue.Supi,
		ServiceType: api.ServiceType_Signalling,
	})
	if err != nil {
		fmt.Printf("ServiceRequest failed: %+v (supi: %s)\n", err, ue.Supi)
		return
	}

	now := time.Now()
	finishTime := now.Sub(startTime)

	if result.StatusCode == api.StatusCode_ERROR {
		fmt.Printf("service failed %s, (supi: %s)\n", result.GetBody(), ue.Supi)
		return false, now, 0, nil
	} else {
		ue.CmState = simulator_context.CmStateConnected
		return true, now, finishTime, nil
	}
}

func triggerAmfFail(count int) {
	logger.ApiLog.Infof("trigger AMF fail")
	_, err := http.Get(fmt.Sprintf("http://10.10.0.18:31118/fail?count=%d", count))
	if err != nil {
		logger.ApiLog.Errorf("http get: %+v", err)
	}
}

func (s *Simulator) ueRegister(ue *simulator_context.UeContext, apiClient api.APIServiceClient, registrationAttempt int) (
	success bool,
	nowTime time.Time,
	completeTime time.Duration,
	redoTime *time.Duration,
) {
	startTime := time.Now()
	for i := 0; i < registrationAttempt; i++ {
		ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second) // T3510
		defer cancel()
		regResult, err := apiClient.Register(ctx, &api.RegisterRequest{
			Supi:            ue.Supi,
			FollowOnRequest: ue.FollowOnRequest,
			ServingPlmn:     ue.ServingPlmnId,
			CipheringAlg:    ue.CipheringAlgStr,
			IntegrityAlg:    ue.IntegrityAlgStr,
			AuthMethod:      ue.AuthData.AuthMethod,
			K:               ue.AuthData.K,
			Opc:             ue.AuthData.Opc,
			Op:              ue.AuthData.Op,
			Amf:             ue.AuthData.AMF,
			Sqn:             ue.AuthData.SQN,
		})
		if err != nil {
			fmt.Printf("Registration failed: %+v (supi: %s)\n", err, ue.Supi)
			// return false, time.Now(), 0, nil
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			continue
		}

		now := time.Now()
		finishTime := now.Sub(startTime)

		if regResult.GetStatusCode() == api.StatusCode_ERROR {
			fmt.Printf("Registration failed: %s (supi: %s)\n", regResult.GetBody(), ue.Supi)
			return false, now, 0, nil
		} else {
			resultUe := regResult.GetUeContext()
			ue.RmState = resultUe.GetRmState()
			ue.CmState = resultUe.GetCmState()
			ue.AmfUeNgapId = resultUe.GetAmfUeNgapId()
			ue.RanUeNgapId = resultUe.GetRanUeNgapId()
			ue.DLCount = security.Count(resultUe.GetNasDownlinkCount())
			ue.ULCount = security.Count(resultUe.GetNasUplinkCount())
			if regResult.RestartCount != 0 {
				restartTime := time.Unix(0, regResult.RestartTimestamp)
				restartFinishTime := now.Sub(restartTime)
				return true, now, finishTime, &restartFinishTime
			} else {
				// fmt.Printf("%s, %d\n", ue.Supi, finishTime.Milliseconds())
				return true, now, finishTime, nil
			}
		}
	}
	return false, time.Now(), 0, nil
}

func (s *Simulator) ueDeregister(ue *simulator_context.UeContext, apiClient api.APIServiceClient) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	startTime := time.Now()
	result, err := apiClient.Deregister(ctx, &api.DeregisterRequest{Supi: ue.Supi})
	if err != nil {
		fmt.Printf("Deregistration failed: %+v (supi: %s)\n", err, ue.Supi)
		return
	}
	finishTime := time.Since(startTime)

	if result.GetStatusCode() == api.StatusCode_ERROR {
		fmt.Printf("Deregistration failed: %s (supi: %s)\n", result.GetBody(), ue.Supi)
	} else {
		fmt.Printf("Deregistration success (supi: %s, expand %+v)\n", ue.Supi, finishTime)
		ue.CmState = simulator_context.CmStateIdle
		ue.RmState = simulator_context.RmStateDeregitered
		s.updateUE(ue)
	}
}

func (s *Simulator) AllUeDeregister(ranName string) {
	apiClient, err := s.connectToRan(ranName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// find all UE and deregister
	cur, err := s.dbClient.Database().Collection("ue").Find(context.TODO(), bson.M{"servingRan": ranName})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No UE found in DB")
		} else {
			fmt.Printf("Find UE error: %+v\n", err)
		}
		return
	}

	wg := sync.WaitGroup{}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		ue := new(simulator_context.UeContext)
		err := cur.Decode(ue)
		if err != nil {
			fmt.Printf("decode ue error: %+v\n", err)
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			s.ueDeregister(ue, apiClient)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

func (s *Simulator) SingleUeDeregister(supi string) {
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

	s.ueDeregister(ue, apiClient)
}

func (s *Simulator) ConnectToAMF(address string, ranName string) {
	apiClient, err := s.connectToRan(ranName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := apiClient.ConnectAMF(context.TODO(), &api.ConnectAMFRequest{Address: address}); err != nil {
		fmt.Printf("connect AMF error: %+v", err)
	}
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
			SubscribedUeAmbr: &models.AmbrRm{
				Uplink:   ue.UeAmbr.UpLink,
				Downlink: ue.UeAmbr.DownLink,
			},
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
