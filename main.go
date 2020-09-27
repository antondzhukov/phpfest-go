package main

import (
	"flag"
	pb "github.com/antondzhukov/phpfest-go/phpfestproto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math"
	"net"
	"os"
	"syscall"
)

type ConfigT struct {
	Help           bool   `default:"false"`
	ConfigFilepath string `default:""`

	Service struct {
		ListenProto       string `default:"tcp"`
		ListenProtoIsUnix bool   `default:"false"`
		ListenAddress     string `default:"localhost:23245"`
	}
}

type AllowedT struct {
	TDouble   [3]float64
	TFloat    [3]float32
	TInt32    [3]int32
	TInt64    [3]int64
	TUint32   [3]uint32
	TUint64   [3]uint64
	TSInt32   [3]int32
	TSInt64   [3]int64
	TFixed32  [3]uint32
	TFixed64  [3]uint64
	TSFixed32 [3]int32
	TSFixed64 [3]int64
	TString   [3]string
	TBytes    [3][]byte
}

var config ConfigT
var allowedTypes AllowedT
var roundRobinCounter int8 = 1

func init() {
	log.SetLevel(log.DebugLevel)

	var help bool

	flag.BoolVar(&help, "help", false, "")
	flag.StringVar(&config.Service.ListenProto, "listen-proto", "unix", "")
	flag.StringVar(&config.Service.ListenAddress, "listen-address", "./phpfest.sock", "")
	flag.Parse()

	config.Help = help

	if config.Service.ListenProto == "unix" || config.Service.ListenProto == "unixpacket" {
		config.Service.ListenProtoIsUnix = true
	}

	if config.Help {
		flag.Usage()
	}
}

func main() {
	allowedTypes = AllowedT{
		TDouble:   [3]float64{math.SmallestNonzeroFloat64, 0.0, math.MaxFloat64},
		TFloat:    [3]float32{math.SmallestNonzeroFloat32, 0.0, math.MaxFloat32},
		TInt32:    [3]int32{math.MinInt32, 0, math.MaxInt32},
		TInt64:    [3]int64{math.MinInt64, 0, math.MaxInt64},
		TUint32:   [3]uint32{0, uint32(math.MaxUint32 / 2), math.MaxUint32},
		TUint64:   [3]uint64{0, uint64(math.MaxUint64 / 2), math.MaxUint64},
		TSInt32:   [3]int32{math.MinInt32, 0, math.MaxInt32},
		TSInt64:   [3]int64{math.MinInt64, 0, math.MaxInt64},
		TFixed32:  [3]uint32{0, uint32(math.MaxUint32 / 2), math.MaxUint32},
		TFixed64:  [3]uint64{0, uint64(math.MaxUint64 / 2), math.MaxUint64},
		TSFixed32: [3]int32{math.MinInt32, 0, math.MaxInt32},
		TSFixed64: [3]int64{math.MinInt64, 0, math.MaxInt64},
		TString:   [3]string{"", "abcdefghijklmn012345", "abcdefghijklmnopqrstuvwxyz0123456789"},
		TBytes:    [3][]byte{[]byte(""), []byte("abcdefghijklmn012345"), []byte("abcdefghijklmnopqrstuvwxyz0123456789")},
	}

	listenAndServe()
}

func listenAndServe() {
	log.Debugf("Try to open connection. Type %s, Address: %s", config.Service.ListenProto, config.Service.ListenAddress)
	if config.Service.ListenProtoIsUnix {
		syscall.Unlink(config.Service.ListenAddress)
	}

	listener, err := net.Listen(config.Service.ListenProto, config.Service.ListenAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	if config.Service.ListenProtoIsUnix {
		os.Chmod(config.Service.ListenAddress, 0666)
	}

	log.Debugf("Listening on proto %s %s", config.Service.ListenProto, config.Service.ListenAddress)

	server := grpc.NewServer()
	pb.RegisterPHPFestService(server, &pb.PHPFestService{GetMessage: getMessage})

	err = server.Serve(listener)
	if err != nil {
		log.Panicf("Error in serve: %s", err.Error())
	}
}

func getMessage(ctx context.Context, in *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	switchRoundRobin()

	if !in.Tbool {
		message := pb.TypesMessage{
			Tbool:     true,
			Tdouble:   allowedTypes.TDouble[roundRobinCounter-1],
			Tfloat:    allowedTypes.TFloat[roundRobinCounter-1],
			Tint32:    allowedTypes.TInt32[roundRobinCounter-1],
			Tint64:    allowedTypes.TInt64[roundRobinCounter-1],
			Tuint32:   allowedTypes.TUint32[roundRobinCounter-1],
			Tuint64:   allowedTypes.TUint64[roundRobinCounter-1],
			Tsint32:   allowedTypes.TSInt32[roundRobinCounter-1],
			Tsint64:   allowedTypes.TSInt64[roundRobinCounter-1],
			Tfixed32:  allowedTypes.TFixed32[roundRobinCounter-1],
			Tfixed64:  allowedTypes.TFixed64[roundRobinCounter-1],
			Tsfixed32: allowedTypes.TSFixed32[roundRobinCounter-1],
			Tsfixed64: allowedTypes.TSFixed64[roundRobinCounter-1],
			Tstring:   allowedTypes.TString[roundRobinCounter-1],
			Tbyte:     allowedTypes.TBytes[roundRobinCounter-1],
		}

		return &pb.GetMessageResponse{
			Tbool:     true,
			Tdouble:   allowedTypes.TDouble[roundRobinCounter-1],
			Tfloat:    allowedTypes.TFloat[roundRobinCounter-1],
			Tint32:    allowedTypes.TInt32[roundRobinCounter-1],
			Tint64:    allowedTypes.TInt64[roundRobinCounter-1],
			Tuint32:   allowedTypes.TUint32[roundRobinCounter-1],
			Tuint64:   allowedTypes.TUint64[roundRobinCounter-1],
			Tsint32:   allowedTypes.TSInt32[roundRobinCounter-1],
			Tsint64:   allowedTypes.TSInt64[roundRobinCounter-1],
			Tfixed32:  allowedTypes.TFixed32[roundRobinCounter-1],
			Tfixed64:  allowedTypes.TFixed64[roundRobinCounter-1],
			Tsfixed32: allowedTypes.TSFixed32[roundRobinCounter-1],
			Tsfixed64: allowedTypes.TSFixed64[roundRobinCounter-1],
			Tstring:   allowedTypes.TString[roundRobinCounter-1],
			Tbyte:     allowedTypes.TBytes[roundRobinCounter-1],
			Tmessage:  &message,
		}, nil
	}

	message := pb.TypesMessage{
		Tbool:     in.Tbool,
		Tdouble:   in.Tdouble,
		Tfloat:    in.Tfloat,
		Tint32:    in.Tint32,
		Tint64:    in.Tint64,
		Tuint32:   in.Tuint32,
		Tuint64:   in.Tuint64,
		Tsint32:   in.Tsint32,
		Tsint64:   in.Tsint64,
		Tfixed32:  in.Tfixed32,
		Tfixed64:  in.Tfixed64,
		Tsfixed32: in.Tsfixed32,
		Tsfixed64: in.Tsfixed64,
		Tstring:   in.Tstring,
		Tbyte:     in.Tbyte,
	}

	return &pb.GetMessageResponse{
		Tbool:     in.Tbool,
		Tdouble:   in.Tdouble,
		Tfloat:    in.Tfloat,
		Tint32:    in.Tint32,
		Tint64:    in.Tint64,
		Tuint32:   in.Tuint32,
		Tuint64:   in.Tuint64,
		Tsint32:   in.Tsint32,
		Tsint64:   in.Tsint64,
		Tfixed32:  in.Tfixed32,
		Tfixed64:  in.Tfixed64,
		Tsfixed32: in.Tsfixed32,
		Tsfixed64: in.Tsfixed64,
		Tstring:   in.Tstring,
		Tbyte:     in.Tbyte,
		Tmessage:  &message,
	}, nil
}

func switchRoundRobin() {
	roundRobinCounter++
	if roundRobinCounter > 3 {
		roundRobinCounter = 1
	}
}
