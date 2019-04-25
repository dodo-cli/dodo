package command

import (
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/spf13/cobra"
)

type options struct {
	file        string
	list        bool
	interactive bool
	remove      bool
	noRemove    bool
	build       bool
	noCache     bool
	pull        bool
	user        string
	workdir     string
	volumes     []string
	volumesFrom []string
	environment []string
	publish     []string
}

func (opts *options) createFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.SetInterspersed(false)

	flags.StringVarP(
		&opts.file, "file", "f", "",
		"specify a dodo configuration file")
	flags.BoolVarP(
		&opts.list, "list", "", false,
		"list all available backdrop configurations")
	flags.BoolVarP(
		&opts.interactive, "interactive", "i", false,
		"run an interactive session")
	flags.BoolVarP(
		&opts.remove, "rm", "", false,
		"automatically remove the container when it exits")
	flags.BoolVarP(
		&opts.noRemove, "no-rm", "", false,
		"keep the container after it exits")
	flags.BoolVarP(
		&opts.build, "build", "", false,
		"always build an image, even if already exists")
	flags.BoolVarP(
		&opts.noCache, "no-cache", "", false,
		"do not use cache when building the image")
	flags.BoolVarP(
		&opts.pull, "pull", "", false,
		"always attempt to pull a newer version of the image")
	flags.StringVarP(
		&opts.user, "user", "u", "",
		"Username or UID (format: <name|uid>[:<group|gid>])")
	flags.StringVarP(
		&opts.workdir, "workdir", "w", "",
		"working directory inside the container")
	flags.StringArrayVarP(
		&opts.volumes, "volume", "v", []string{},
		"Bind mount a volume")
	flags.StringArrayVarP(
		&opts.volumesFrom, "volumes-from", "", []string{},
		"Mount volumes from the specified container(s)")
	flags.StringArrayVarP(
		&opts.environment, "env", "e", []string{},
		"Set environment variables")
	flags.StringArrayVarP(
		&opts.publish, "publish", "p", []string{},
		"Publish a container's port(s) to the host")
}

func (opts *options) createConfig(command []string) (*types.Backdrop, error) {
	config := &types.Backdrop{
		Image: &types.Image{
			ForceRebuild: opts.build,
			NoCache:      opts.noCache,
			ForcePull:    opts.pull,
		},
		Interactive: opts.interactive,
		User:        opts.user,
		WorkingDir:  opts.workdir,
		VolumesFrom: opts.volumesFrom,
		Command:     command,
	}

	if opts.noRemove {
		remove := false
		config.Remove = &remove
	}
	if opts.remove {
		remove := true
		config.Remove = &remove
	}

	for _, volume := range opts.volumes {
		decoded, err := types.DecodeVolume("cli", volume)
		if err != nil {
			return nil, err
		}
		config.Volumes = append(config.Volumes, decoded)
	}

	for _, env := range opts.environment {
		decoded, err := types.DecodeKeyValue("cli", env)
		if err != nil {
			return nil, err
		}
		config.Environment = append(config.Environment, decoded)
	}

	for _, port := range opts.publish {
		decoded, err := types.DecodePort("cli", port)
		if err != nil {
			return nil, err
		}
		config.Ports = append(config.Ports, decoded)
	}

	return config, nil
}