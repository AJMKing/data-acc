package ansible

import (
	"fmt"
	"github.com/RSE-Cambridge/data-acc/internal/pkg/registry"
	"log"
	"os/exec"
	"path"
)

func getMountDir(volume registry.Volume) string {
	// TODO: what about the environment variables that are being set? should share logic with here
	if volume.MultiJob {
		return fmt.Sprintf("/dac/multi_job_buffer/%s", volume.UUID)
	}
	return fmt.Sprintf("/dac/job_buffer/%s", volume.UUID)
}

func mount(fsType FSType, volume registry.Volume, brickAllocations []registry.BrickAllocation) error {
	var primaryBrickHost string
	for _, allocation := range brickAllocations {
		if allocation.AllocatedIndex == 0 {
			primaryBrickHost = allocation.Hostname
			break
		}
	}

	if primaryBrickHost == "" {
		log.Panicf("failed to find primary brick for volume: %s", volume.Name)
	}

	if fsType == BeegFS {
		// Write out the config needed, and do the mount using ansible
		// TODO: Move Lustre mount here that is done below
		executeAnsibleMount(fsType, volume, brickAllocations)
	}

	var mountDir = getMountDir(volume)
	for _, attachment := range volume.Attachments {

		if err := mkdir(attachment.Hostname, mountDir); err != nil {
			return err
		}
		if err := mountRemoteFilesystem(fsType, attachment.Hostname,
			primaryBrickHost, volume.UUID, mountDir); err != nil {
			return err
		}

		if !volume.MultiJob && volume.AttachAsSwapBytes > 0 {
			swapDir := path.Join(mountDir, "/swap")
			if err := mkdir(attachment.Hostname, swapDir); err != nil {
				return err
			}

			// TODO: swapmb := int(volume.AttachAsSwapBytes/1024)
			swapmb := 2
			swapFile := path.Join(swapDir, fmt.Sprintf("/%s", attachment.Hostname))
			loopback := fmt.Sprintf("/dev/loop%d", volume.ClientPort)
			if err := createSwap(attachment.Hostname, swapmb, swapFile, loopback); err != nil {
				return err
			}

			if err := swapOn(attachment.Hostname, loopback); err != nil {
				return err
			}
		}

		if !volume.MultiJob && volume.AttachPrivateNamespace {
			privateDir := path.Join(mountDir, fmt.Sprintf("/private/%s", attachment.Hostname))
			if err := mkdir(attachment.Hostname, privateDir); err != nil {
				return err
			}
			chown(attachment.Hostname, volume.Owner, volume.Group, privateDir)
		}

		sharedDir := path.Join(mountDir, "/shared")
		if err := mkdir(attachment.Hostname, sharedDir); err != nil {
			return err
		}
		chown(attachment.Hostname, volume.Owner, volume.Group, sharedDir)
	}
	// TODO on error should we always call umount? maybe?
	// TODO move to ansible style automation or preamble?
	return nil
}

func umount(fsType FSType, volume registry.Volume, brickAllocations []registry.BrickAllocation) error {
	log.Println("FAKE Umount for:", volume.Name)
	var mountDir = getMountDir(volume)

	for _, attachment := range volume.Attachments {
		if !volume.MultiJob && volume.AttachAsSwapBytes > 0 {
			swapFile := path.Join(mountDir, fmt.Sprintf("/swap/%s", attachment.Hostname)) // TODO share?
			loopback := fmt.Sprintf("/dev/loop%d", volume.ClientPort)                     // TODO share?
			if err := swapOff(attachment.Hostname, loopback); err != nil {
				return err
			}
			if err := detachLoopback(attachment.Hostname, loopback); err != nil {
				return err
			}
			if err := removeSubtree(attachment.Hostname, swapFile); err != nil {
				return err
			}
		}
		if fsType == Lustre {
			if err := umountLustre(attachment.Hostname, mountDir); err != nil {
				return err
			}
		}
		if err := removeSubtree(attachment.Hostname, mountDir); err != nil {
			return err
		}
	}

	if fsType == BeegFS {
		// TODO: Move Lustre unmount here that is done below
		executeAnsibleUnmount(fsType, volume, brickAllocations)
		// TODO: this makes copy out much harder in its current form :(
	}

	return nil
}

func createSwap(hostname string, swapMb int, filename string, loopback string) error {
	file := fmt.Sprintf("dd if=/dev/zero of=%s bs=1024 count=%d && sudo chmod 0600 %s",
		filename, swapMb*1024, filename)
	if err := remoteExecuteCmd(hostname, file); err != nil {
		return err
	}
	device := fmt.Sprintf("losetup %s %s", loopback, filename)
	if err := remoteExecuteCmd(hostname, device); err != nil {
		return err
	}
	swap := fmt.Sprintf("mkswap %s", loopback)
	return remoteExecuteCmd(hostname, swap)
}

func swapOn(hostname string, loopback string) error {
	return remoteExecuteCmd(hostname, fmt.Sprintf("swapon %s", loopback))
}

func swapOff(hostname string, loopback string) error {
	return remoteExecuteCmd(hostname, fmt.Sprintf("swapoff %s", loopback))
}

func detachLoopback(hostname string, loopback string) error {
	return remoteExecuteCmd(hostname, fmt.Sprintf("losetup -d %s", loopback))
}

func chown(hostname string, owner uint, group uint, directory string) error {
	return remoteExecuteCmd(hostname, fmt.Sprintf("chown %d:%d %s", owner, group, directory))
}

func umountLustre(hostname string, directory string) error {
	return remoteExecuteCmd(hostname, fmt.Sprintf("umount -l %s", directory))
}

func removeSubtree(hostname string, directory string) error {
	return remoteExecuteCmd(hostname, fmt.Sprintf("rm -rf %s", directory))
}

func mountRemoteFilesystem(fsType FSType, hostname string, mgtHost string, fsname string, directory string) error {
	if fsType == Lustre {
		return mountLustre(hostname, mgtHost, fsname, directory)
	} else if fsType == BeegFS {
		return mountBeegFS(hostname, mgtHost, fsname, directory)
	}
	return fmt.Errorf("mount unsuported by filesystem type %s", fsType)
}

func mountLustre(hostname string, mgtHost string, fsname string, directory string) error {
	if err := remoteExecuteCmd(hostname, "modprobe -v lustre"); err != nil {
		return err
	}
	return remoteExecuteCmd(hostname, fmt.Sprintf(
		"mount -t lustre %s-opa@o2ib1:/%s %s", mgtHost, fsname, directory)) //TODO: make the lustre nid configurable this is filth
}

func mountBeegFS(hostname string, mgtHost string, fsname string, directory string) error {
	// Ansible mounts beegfs at /dac/beegfs/<fsname>, link into above location here
	// First remove the directory, then replace with a symbolic link
	if err := removeSubtree(hostname, directory); err != nil {
		return err
	}
	return remoteExecuteCmd(hostname, fmt.Sprintf("ln -s /dac/beegfs/%s %s", fsname, directory))
}

func mkdir(hostname string, directory string) error {
	return remoteExecuteCmd(hostname, fmt.Sprintf("mkdir -p %s", directory))
}

func remoteExecuteCmd(hostname string, cmdStr string) error {
	log.Println("SSH to:", hostname, "with command:", cmdStr)

	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null", hostname, "sudo", cmdStr)
	output, err := cmd.CombinedOutput()

	if err == nil {
		log.Println("Completed remote ssh run:", cmdStr)
		log.Println(string(output))
		return nil
	} else {
		log.Println("Error in remove ssh run:", string(output))
		return err
	}
}
