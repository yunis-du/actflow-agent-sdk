package sdk

import (
	"context"
	"net"

	"github.com/yunis-du/actflow-agent-sdk/go-sdk/pb"
	"google.golang.org/grpc"
)

// Server wraps an Agent implementation as a gRPC server.
type Server struct {
	agent  Agent
	server *grpc.Server
}

// NewServer creates a new agent server.
func NewServer(agent Agent) *Server {
	return &Server{
		agent: agent,
	}
}

// Serve starts the gRPC server on the given address.
func (s *Server) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()
	pb.RegisterAgentServiceServer(s.server, &agentServiceImpl{agent: s.agent})

	return s.server.Serve(lis)
}

// GracefulStop gracefully stops the server.
func (s *Server) GracefulStop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// Stop immediately stops the server.
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
}

// agentServiceImpl implements the gRPC AgentService.
type agentServiceImpl struct {
	pb.UnimplementedAgentServiceServer
	agent Agent
}

// Run implements the Run RPC method.
func (s *agentServiceImpl) Run(req *pb.RunRequest, stream pb.AgentService_RunServer) error {
	// Convert proto context to SDK context
	execCtx := protoContextToContext(req.Ctx)

	// Convert inputs
	inputs := protoValueToAny(req.Inputs)

	// Create log channel
	logCh := make(chan string, 1024)

	// Run agent in goroutine
	outputCh := make(chan *Output, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(logCh)
		output, err := s.agent.Run(stream.Context(), req.Nid, execCtx, inputs, logCh)
		if err != nil {
			errCh <- err
			return
		}
		outputCh <- output
	}()

	// Stream logs
	for log := range logCh {
		update := &pb.AgentUpdate{
			RelayMessage: &pb.AgentUpdate_Log{Log: log},
		}
		if err := stream.Send(update); err != nil {
			return err
		}
	}

	// Wait for output or error
	select {
	case err := <-errCh:
		// Send error as failed output
		update := &pb.AgentUpdate{
			RelayMessage: &pb.AgentUpdate_Output{
				Output: &pb.AgentOutput{
					Status: pb.NodeExecutionStatus_FAILED,
					Error:  err.Error(),
				},
			},
		}
		return stream.Send(update)

	case output := <-outputCh:
		// Send output
		update := &pb.AgentUpdate{
			RelayMessage: &pb.AgentUpdate_Output{
				Output: outputToProto(output),
			},
		}
		return stream.Send(update)

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// Shutdown implements the Shutdown RPC method.
func (s *agentServiceImpl) Shutdown(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	if err := s.agent.Shutdown(ctx); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
