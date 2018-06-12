package main

func install(u string) error {
	i, err := newInstaller(u)

	if err != nil {
		return err
	}

	err = i.InstallModule()

	if err != nil {
		return err
	}

	return i.InstallCommands()
}
