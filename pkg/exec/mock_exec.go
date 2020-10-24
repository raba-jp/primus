package exec

import (
	"context"
	"syscall"
	"io"
	mock "github.com/stretchr/testify/mock"
)

type MockInterface struct {
	mock.Mock
}

type InterfaceCommandArgs struct {
	Cmd          string
	CmdAnything  bool
	Args         []string
	ArgsAnything bool
}

type InterfaceCommandReturns struct {
	Cmd func() Cmd
}

type InterfaceCommandExpectation struct {
	Args    InterfaceCommandArgs
	Returns InterfaceCommandReturns
}

func (_m *MockInterface) ApplyCommandExpectation(e InterfaceCommandExpectation) {
	var args []interface{}
	if e.Args.CmdAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Cmd)
	}
	if e.Args.ArgsAnything {
		args = append(args, mock.Anything)
	} else if e.Args.Args != nil {
		for _, arg := range e.Args.Args {
			args = append(args, arg)
		}
	}
	_m.On("Command", args...).Return(e.Returns.Cmd())
}

func (_m *MockInterface) ApplyCommandExpectations(expectations []InterfaceCommandExpectation) {
	for _, e := range expectations {
		_m.ApplyCommandExpectation(e)
	}
}

// Command provides a mock function with given fields: cmd, args
func (_m *MockInterface) Command(cmd string, args ...string) Cmd {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, cmd)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var c Cmd
	if rf, ok := ret.Get(0).(func(string, ...string) Cmd); ok {
		c = rf(cmd, args...)
	} else {
		if ret.Get(0) != nil {
			c = ret.Get(0).(Cmd)
		}
	}

	return c
}

type InterfaceCommandContextArgs struct {
	Ctx          context.Context
	CtxAnything  bool
	Cmd          string
	CmdAnything  bool
	Args         []string
	ArgsAnything bool
}

type InterfaceCommandContextReturns struct {
	Cmd func()Cmd
}

type InterfaceCommandContextExpectation struct {
	Args    InterfaceCommandContextArgs
	Returns InterfaceCommandContextReturns
}

func (_m *MockInterface) ApplyCommandContextExpectation(e InterfaceCommandContextExpectation) {
	var args []interface{}
	if e.Args.CtxAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Ctx)
	}
	if e.Args.CmdAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Cmd)
	}
	if e.Args.ArgsAnything {
		args = append(args, mock.Anything)
	} else if e.Args.Args != nil {
		for _, arg := range e.Args.Args {
			args = append(args, arg)
		}
	}
	_m.On("CommandContext", args...).Return(e.Returns.Cmd())
}

func (_m *MockInterface) ApplyCommandContextExpectations(expectations []InterfaceCommandContextExpectation) {
	for _, e := range expectations {
		_m.ApplyCommandContextExpectation(e)
	}
}

func (_m *MockInterface) CommandContext(ctx context.Context, cmd string, args ...string) Cmd {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, cmd)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var c Cmd
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) Cmd); ok {
		c = rf(ctx, cmd, args...)
	} else {
		if ret.Get(0) != nil {
			c = ret.Get(0).(Cmd)
		}
	}

	return c
}

type InterfaceLookPathArgs struct {
	File         string
	FileAnything bool
}

type InterfaceLookPathReturns struct {
	Path string
	Err error
}

type InterfaceLookPathExpectation struct {
	Args    InterfaceLookPathArgs
	Returns InterfaceLookPathReturns
}

func (_m *MockInterface) ApplyLookPathExpectation(e InterfaceLookPathExpectation) {
	var args []interface{}
	if e.Args.FileAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.File)
	}
	_m.On("LookPath", args...).Return(e.Returns.Path, e.Returns.Err)
}

func (_m *MockInterface) ApplyLookPathExpectations(expectations []InterfaceLookPathExpectation) {
	for _, e := range expectations {
		_m.ApplyLookPathExpectation(e)
	}
}

// LookPath provides a mock function with given fields: file
func (_m *MockInterface) LookPath(file string) (string, error) {
	ret := _m.Called(file)

	var path string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		path = rf(file)
	} else {
		path = ret.Get(0).(string)
	}

	var err error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		err = rf(file)
	} else {
		err = ret.Error(1)
	}

	return path, err
}

type MockCmd struct {
	mock.Mock
}

type CmdCombinedOutputReturns struct {
	Output []byte
	Err error
}

type CmdCombinedOutputExpectation struct {
	Returns CmdCombinedOutputReturns
}

func (_m *MockCmd) ApplyCombinedOutputExpectation(e CmdCombinedOutputExpectation) {
	var args []interface{}
	_m.On("CombinedOutput", args...).Return(e.Returns.Output, e.Returns.Err)
}

func (_m *MockCmd) ApplyCombinedOutputExpectations(expectations []CmdCombinedOutputExpectation) {
	for _, e := range expectations {
		_m.ApplyCombinedOutputExpectation(e)
	}
}

func (_m *MockCmd) CombinedOutput() ([]byte, error) {
	ret := _m.Called()

	var output []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		output = rf()
	} else {
		if ret.Get(0) != nil {
			output = ret.Get(0).([]byte)
		}
	}

	var err error
	if rf, ok := ret.Get(1).(func() error); ok {
		err = rf()
	} else {
		err = ret.Error(1)
	}

	return output, err
}

type CmdOutputReturns struct {
	Output []byte
	Err error
}

type CmdOutputExpectation struct {
	Returns CmdOutputReturns
}

func (_m *MockCmd) ApplyOutputExpectation(e CmdOutputExpectation) {
	var args []interface{}
	_m.On("Output", args...).Return(e.Returns.Output, e.Returns.Err)
}

func (_m *MockCmd) ApplyOutputExpectations(expectations []CmdOutputExpectation) {
	for _, e := range expectations {
		_m.ApplyOutputExpectation(e)
	}
}

func (_m *MockCmd) Output() ([]byte, error) {
	ret := _m.Called()

	var output []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		output = rf()
	} else {
		if ret.Get(0) != nil {
			output = ret.Get(0).([]byte)
		}
	}

	var err error
	if rf, ok := ret.Get(1).(func() error); ok {
		err = rf()
	} else {
		err = ret.Error(1)
	}

	return output, err
}

type CmdRunReturns struct {
	Err error
}

type CmdRunExpectation struct {
	Returns CmdRunReturns
}

func (_m *MockCmd) ApplyRunExpectation(e CmdRunExpectation) {
	var args []interface{}
	_m.On("Run", args...).Return(e.Returns.Err)
}

func (_m *MockCmd) ApplyRunExpectations(expectations []CmdRunExpectation) {
	for _, e := range expectations {
		_m.ApplyRunExpectation(e)
	}
}

func (_m *MockCmd) Run() error {
	ret := _m.Called()

	var err error
	if rf, ok := ret.Get(0).(func() error); ok {
		err = rf()
	} else {
		err = ret.Error(0)
	}

	return err
}

type CmdSetDirArgs struct {
	Dir         string
	DirAnything bool
}

type CmdSetDirExpectation struct {
	Args CmdSetDirArgs
}

func (_m *MockCmd) ApplySetDirExpectation(e CmdSetDirExpectation) {
	var args []interface{}
	if e.Args.DirAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Dir)
	}
	_m.On("SetDir", args...)
}

func (_m *MockCmd) ApplySetDirExpectations(expectations []CmdSetDirExpectation) {
	for _, e := range expectations {
		_m.ApplySetDirExpectation(e)
	}
}

func (_m *MockCmd) SetDir(dir string) {
	_m.Called(dir)
}

type CmdSetEnvArgs struct {
	Env         []string
	EnvAnything bool
}

type CmdSetEnvExpectation struct {
	Args CmdSetEnvArgs
}

func (_m *MockCmd) ApplySetEnvExpectation(e CmdSetEnvExpectation) {
	var args []interface{}
	if e.Args.EnvAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Env)
	}
	_m.On("SetEnv", args...)
}

func (_m *MockCmd) ApplySetEnvExpectations(expectations []CmdSetEnvExpectation) {
	for _, e := range expectations {
		_m.ApplySetEnvExpectation(e)
	}
}

func (_m *MockCmd) SetEnv(env []string) {
	_m.Called(env)
}

type CmdSetStderrArgs struct {
	Out         io.Writer
	OutAnything bool
}

type CmdSetStderrExpectation struct {
	Args CmdSetStderrArgs
}

func (_m *MockCmd) ApplySetStderrExpectation(e CmdSetStderrExpectation) {
	var args []interface{}
	if e.Args.OutAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Out)
	}
	_m.On("SetStderr", args...)
}

func (_m *MockCmd) ApplySetStderrExpectations(expectations []CmdSetStderrExpectation) {
	for _, e := range expectations {
		_m.ApplySetStderrExpectation(e)
	}
}

func (_m *MockCmd) SetStderr(out io.Writer) {
	_m.Called(out)
}

type CmdSetStdinArgs struct {
	In         io.Reader
	InAnything bool
}

type CmdSetStdinExpectation struct {
	Args CmdSetStdinArgs
}

func (_m *MockCmd) ApplySetStdinExpectation(e CmdSetStdinExpectation) {
	var args []interface{}
	if e.Args.InAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.In)
	}
	_m.On("SetStdin", args...)
}

func (_m *MockCmd) ApplySetStdinExpectations(expectations []CmdSetStdinExpectation) {
	for _, e := range expectations {
		_m.ApplySetStdinExpectation(e)
	}
}

func (_m *MockCmd) SetStdin(in io.Reader) {
	_m.Called(in)
}

type CmdSetStdoutArgs struct {
	Out         io.Writer
	OutAnything bool
}

type CmdSetStdoutExpectation struct {
	Args CmdSetStdoutArgs
}

func (_m *MockCmd) ApplySetStdoutExpectation(e CmdSetStdoutExpectation) {
	var args []interface{}
	if e.Args.OutAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Out)
	}
	_m.On("SetStdout", args...)
}

func (_m *MockCmd) ApplySetStdoutExpectations(expectations []CmdSetStdoutExpectation) {
	for _, e := range expectations {
		_m.ApplySetStdoutExpectation(e)
	}
}

func (_m *MockCmd) SetStdout(out io.Writer) {
	_m.Called(out)
}

type CmdSetSysProcAttrArgs struct {
	Proc         *syscall.SysProcAttr
	ProcAnything bool
}

type CmdSetSysProcAttrExpectation struct {
	Args CmdSetSysProcAttrArgs
}

func (_m *MockCmd) ApplySetSysProcAttrExpectation(e CmdSetSysProcAttrExpectation) {
	var args []interface{}
	if e.Args.ProcAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Proc)
	}
	_m.On("SetSysProcAttr", args...)
}

func (_m *MockCmd) ApplySetSysProcAttrExpectations(expectations []CmdSetSysProcAttrExpectation) {
	for _, e := range expectations {
		_m.ApplySetSysProcAttrExpectation(e)
	}
}

func (_m *MockCmd) SetSysProcAttr(proc *syscall.SysProcAttr) {
	_m.Called(proc)
}

type CmdStartReturns struct {
	Err error
}

type CmdStartExpectation struct {
	Returns CmdStartReturns
}

func (_m *MockCmd) ApplyStartExpectation(e CmdStartExpectation) {
	var args []interface{}
	_m.On("Start", args...).Return(e.Returns.Err)
}

func (_m *MockCmd) ApplyStartExpectations(expectations []CmdStartExpectation) {
	for _, e := range expectations {
		_m.ApplyStartExpectation(e)
	}
}

func (_m *MockCmd) Start() error {
	ret := _m.Called()

	var err error
	if rf, ok := ret.Get(0).(func() error); ok {
		err = rf()
	} else {
		err = ret.Error(0)
	}

	return err
}

type CmdStderrPipeReturns struct {
	Output io.ReadCloser
	Err error
}

type CmdStderrPipeExpectation struct {
	Returns CmdStderrPipeReturns
}

func (_m *MockCmd) ApplyStderrPipeExpectation(e CmdStderrPipeExpectation) {
	var args []interface{}
	_m.On("StderrPipe", args...).Return(e.Returns.Output, e.Returns.Err)
}

func (_m *MockCmd) ApplyStderrPipeExpectations(expectations []CmdStderrPipeExpectation) {
	for _, e := range expectations {
		_m.ApplyStderrPipeExpectation(e)
	}
}

// StderrPipe provides a mock function with given fields:
func (_m *MockCmd) StderrPipe() (io.ReadCloser, error) {
	ret := _m.Called()

	var output io.ReadCloser
	if rf, ok := ret.Get(0).(func() io.ReadCloser); ok {
		output = rf()
	} else {
		if ret.Get(0) != nil {
			output = ret.Get(0).(io.ReadCloser)
		}
	}

	var err error
	if rf, ok := ret.Get(1).(func() error); ok {
		err = rf()
	} else {
		err = ret.Error(1)
	}

	return output, err
}

type CmdStdoutPipeReturns struct {
	Output io.ReadCloser
	Err error
}

type CmdStdoutPipeExpectation struct {
	Returns CmdStdoutPipeReturns
}

func (_m *MockCmd) ApplyStdoutPipeExpectation(e CmdStdoutPipeExpectation) {
	var args []interface{}
	_m.On("StdoutPipe", args...).Return(e.Returns.Output, e.Returns.Err)
}

func (_m *MockCmd) ApplyStdoutPipeExpectations(expectations []CmdStdoutPipeExpectation) {
	for _, e := range expectations {
		_m.ApplyStdoutPipeExpectation(e)
	}
}

// StdoutPipe provides a mock function with given fields:
func (_m *MockCmd) StdoutPipe() (io.ReadCloser, error) {
	ret := _m.Called()

	var output io.ReadCloser
	if rf, ok := ret.Get(0).(func() io.ReadCloser); ok {
		output = rf()
	} else {
		if ret.Get(0) != nil {
			output = ret.Get(0).(io.ReadCloser)
		}
	}

	var err error
	if rf, ok := ret.Get(1).(func() error); ok {
		err = rf()
	} else {
		err = ret.Error(1)
	}

	return output, err
}

type CmdStopExpectation struct {
}

func (_m *MockCmd) ApplyStopExpectation(e CmdStopExpectation) {
	var args []interface{}
	_m.On("Stop", args...)
}

func (_m *MockCmd) ApplyStopExpectations(expectations []CmdStopExpectation) {
	for _, e := range expectations {
		_m.ApplyStopExpectation(e)
	}
}

// Stop provides a mock function with given fields:
func (_m *MockCmd) Stop() {
	_m.Called()
}

type CmdWaitReturns struct {
	Err error
}

type CmdWaitExpectation struct {
	Returns CmdWaitReturns
}

func (_m *MockCmd) ApplyWaitExpectation(e CmdWaitExpectation) {
	var args []interface{}
	_m.On("Wait", args...).Return(e.Returns.Err)
}

func (_m *MockCmd) ApplyWaitExpectations(expectations []CmdWaitExpectation) {
	for _, e := range expectations {
		_m.ApplyWaitExpectation(e)
	}
}

// Wait provides a mock function with given fields:
func (_m *MockCmd) Wait() error {
	ret := _m.Called()

	var err error
	if rf, ok := ret.Get(0).(func() error); ok {
		err = rf()
	} else {
		err = ret.Error(0)
	}

	return err
}
