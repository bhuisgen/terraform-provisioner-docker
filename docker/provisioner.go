package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/armon/circbuf"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/go-linereader"
)

const (
	maxBufSize = 8 * 1024
)

func Provisioner() terraform.ResourceProvisioner {
	return &schema.Provisioner{
		Schema: map[string]*schema.Schema{
			"container": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"command": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"environment": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},

		ApplyFunc: applyFn,
	}
}

func applyFn(ctx context.Context) error {
	data := ctx.Value(schema.ProvConfigDataKey).(*schema.ResourceData)
	o := ctx.Value(schema.ProvOutputKey).(terraform.UIOutput)

	v, ok := data.GetOk("command")
	if !ok {
		return fmt.Errorf("docker provisioner command must be a non-empty string")
	}

	command := stringListToStringSlice(v.([]interface{}))
	for _, v := range command {
		if v == "" {
			return fmt.Errorf("values for command may not be empty")
		}
	}

	environment := data.Get("environment").(map[string]interface{})

	var env []string
	for k := range environment {
		entry := fmt.Sprintf("%s=%s", k, environment[k].(string))
		env = append(env, entry)
	}

	var cmdargs []string
	cmdargs = []string{"docker", "exec", data.Get("container").(string)}
	for _, v := range command {
		cmdargs = append(cmdargs, v)
	}

	pr, pw, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to initialize pipe for output: %s", err)
	}

	var cmdEnv []string
	cmdEnv = os.Environ()
	cmdEnv = append(cmdEnv, env...)

	cmd := exec.CommandContext(ctx, cmdargs[0], cmdargs[1:]...)
	cmd.Stderr = pw
	cmd.Stdout = pw
	cmd.Env = cmdEnv

	output, _ := circbuf.NewBuffer(maxBufSize)

	tee := io.TeeReader(pr, output)

	copyDoneCh := make(chan struct{})
	go copyOutput(o, tee, copyDoneCh)

	o.Output(fmt.Sprintf("Executing: %q", cmdargs))

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	pw.Close()

	select {
	case <-copyDoneCh:
	case <-ctx.Done():
	}

	if err != nil {
		return fmt.Errorf("Error running command '%s': %v. Output: %s",
			command, err, output.Bytes())
	}

	return nil
}

func copyOutput(o terraform.UIOutput, r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)
	lr := linereader.New(r)
	for line := range lr.Ch {
		o.Output(line)
	}
}

func stringListToStringSlice(stringList []interface{}) []string {
	ret := []string{}
	for _, v := range stringList {
		if v == nil {
			ret = append(ret, "")
			continue
		}
		ret = append(ret, v.(string))
	}
	return ret
}
