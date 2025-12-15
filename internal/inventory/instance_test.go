package inventory

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewInstance verifies that NewInstance returns a correctly initialized instance
func TestNewInstance(t *testing.T) {
	t.Run("New Instance", func(t *testing.T) { testNewInstance_Correct(t) })
	t.Run("New Instance without InstanceID", func(t *testing.T) { testNewInstance_WithoutInstanceID(t) })
}

func testNewInstance_Correct(t *testing.T) {
	id := "id-012345X"
	name := "test-instance"
	provider := AWSProvider
	instanceType := "vm.small"
	az := "north-eu-1"
	status := Running
	tags := []Tag{}
	tz := time.Now()

	instance, err := NewInstance(id, name, provider, instanceType, az, status, tags, tz)

	// Basic check
	assert.Nil(t, err)
	assert.NotNil(t, instance)

	// Parameters Check
	assert.Equal(t, instance.InstanceID, id)
	assert.Equal(t, instance.InstanceName, name)
	assert.Equal(t, instance.Provider, provider)
	assert.Equal(t, instance.InstanceType, instanceType)
	assert.Equal(t, instance.AvailabilityZone, az)
	assert.Equal(t, instance.Status, status)
	assert.Equal(t, instance.ClusterID, "")
	assert.NotZero(t, instance.LastScanTS)
	assert.Equal(t, instance.CreatedAt, tz)
	assert.Equal(t, instance.Age, 1)
	assert.Equal(t, instance.Tags, tags)
	assert.NotNil(t, instance.Expenses)
}

func testNewInstance_WithoutInstanceID(t *testing.T) {
	id := ""
	name := "test-instance"
	provider := AWSProvider
	instanceType := "vm.small"
	az := "north-eu-1"
	status := Running
	tags := []Tag{}
	tz := time.Now()

	instance, err := NewInstance(id, name, provider, instanceType, az, status, tags, tz)

	// Basic check
	assert.Nil(t, instance)
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrMissingInstanceIDCreation.Error())
}

// TestAddTag tests the tag adding operation to an Instance
func TestAddTag(t *testing.T) {
	t.Run("Adding Tag", func(t *testing.T) { testAddTag_Correct(t) })
	t.Run("Adding Keyless Tag", func(t *testing.T) { testAddTag_WithoutKey(t) })
}

func testAddTag_Correct(t *testing.T) {
	i := Instance{}
	tag := Tag{Key: "env", Value: "prod"}

	assert.Zero(t, len(i.Tags))

	err := i.AddTag(tag)
	assert.Nil(t, err)

	assert.Equal(t, len(i.Tags), 1)
	assert.Contains(t, i.Tags, tag)
}

func testAddTag_WithoutKey(t *testing.T) {
	i := Instance{}
	tag := Tag{Key: "", Value: "value"}

	assert.Zero(t, len(i.Tags))

	err := i.AddTag(tag)
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrAddingTagWithoutKey.Error())
	assert.Zero(t, len(i.Tags))
}

// TestAddExpense tests the Expense adding operation to an Instance
func TestAddExpense(t *testing.T) {
	t.Run("Adding Expense", func(t *testing.T) { testAddExpense_Correct(t) })
	t.Run("Adding Expense with negative amount", func(t *testing.T) { testAddExpense_WithNegativeAmount(t) })
}

func testAddExpense_Correct(t *testing.T) {
	i := Instance{InstanceID: "id-012345X"}
	expense := Expense{InstanceID: "id-00000X", Amount: 2.4, Date: time.Now()}

	assert.Zero(t, len(i.Expenses))

	err := i.AddExpense(&expense)
	assert.Nil(t, err)

	assert.Equal(t, len(i.Expenses), 1)
	assert.Contains(t, i.Expenses, expense)
}

func testAddExpense_WithNegativeAmount(t *testing.T) {
	i := Instance{InstanceID: "id-012345X"}
	expense := Expense{InstanceID: "id-00000X", Amount: -2.4, Date: time.Now()}

	assert.Zero(t, len(i.Expenses))

	err := i.AddExpense(&expense)
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrAddingExpenseWithWrongAmount.Error())
	assert.Zero(t, len(i.Tags))
}

// TestInstance_String verifies String method returns expected format
func TestInstance_String(t *testing.T) {
	i := Instance{
		InstanceID:       "i-123",
		InstanceName:     "test",
		Provider:         AWSProvider,
		InstanceType:     "t2.micro",
		AvailabilityZone: "us-east-1a",
		Status:           Running,
		ClusterID:        "cluster-x",
		Expenses:         []Expense{{Amount: 5}},
	}

	str := i.String()
	if !(strings.Contains(str, "test") && strings.Contains(str, "AWS") && strings.Contains(str, "t2.micro")) {
		t.Errorf("unexpected output from String(): %s", str)
	}
}

// TestPrintInstance verifies PrintInstance runs without panic
func TestPrintInstance(t *testing.T) {
	i := Instance{InstanceID: "i-456", InstanceName: "node1"}
	i.PrintInstance()
}
