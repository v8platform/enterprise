package enterprise

import (
	"github.com/stretchr/testify/suite"
	"github.com/v8platform/errors"
	"github.com/v8platform/runner"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var pwd, _ = os.Getwd()

func NewFileIB(path string) TempInfobase {

	ib := TempInfobase{
		File: path,
	}

	return ib
}

type TempInfobase struct {
	File string
}

func (ib TempInfobase) Path() string {
	return ib.File
}

func (ib TempInfobase) ConnectionString() string {
	return "/F" + ib.File
}

func (ib TempInfobase) Values() []string {

	return []string{"file=" + ib.File}
}

type TempCreateInfobase struct {
}

func (ib TempCreateInfobase) Command() string {
	return "CREATEINFOBASE"
}

func (ib TempCreateInfobase) Check() error {
	return nil
}
func (ib TempCreateInfobase) Values() []string {

	return []string{}
}

type EnterpriseTestSuite struct {
	suite.Suite
	TempIB string
}

func TestEnterprise(t *testing.T) {
	suite.Run(t, new(EnterpriseTestSuite))
}

func (t *EnterpriseTestSuite) AfterTest(suite, testName string) {
	t.ClearTempInfoBase()
}

func (t *EnterpriseTestSuite) BeforeTest(suite, testName string) {
	t.CreateTempInfoBase()
}

func (t *EnterpriseTestSuite) SetupSuite() {

	ibPath, _ := ioutil.TempDir("", "1c_DB_")
	t.TempIB = ibPath

}

func (t *EnterpriseTestSuite) CreateTempInfoBase() {

	ib := TempInfobase{File: t.TempIB}

	err := runner.Run(ib, TempCreateInfobase{},
		runner.WithTimeout(30))

	t.Require().NoError(err, errors.GetErrorContext(err))

}

func (t *EnterpriseTestSuite) ClearTempInfoBase() {

	err := os.RemoveAll(t.TempIB)
	t.Require().NoError(err, errors.GetErrorContext(err))
}

func (t *EnterpriseTestSuite) TestRunEpf() {

	epf := path.Join(pwd, "tests", "fixtures", "epf", "Test_Close.epf")

	err := runner.Run(NewFileIB(t.TempIB), ExecuteOptions{
		File: epf},
		runner.WithTimeout(30))

	t.Require().NoError(err, errors.GetErrorContext(err))

}

func (t *EnterpriseTestSuite) TestRunWithParam() {

	epf := path.Join(pwd, "tests", "fixtures", "epf", "Test_Close.epf")

	exec := ExecuteOptions{
		File: epf}.WithParams(map[string]string{"Привет": "мир"})

	err := runner.Run(NewFileIB(t.TempIB), exec,
		runner.WithTimeout(30))

	t.Require().NoError(err, errors.GetErrorContext(err))

}
