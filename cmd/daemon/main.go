package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/opencontainers/image-tools/image"
	kataTypes "github.com/openshift/kata-operator/pkg/apis/kataconfiguration/v1alpha1"
	kataClient "github.com/openshift/kata-operator/pkg/generated/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	var kataOperation string
	flag.StringVar(&kataOperation, "operation", "", "Specify kata operations. Valid options are 'install', 'upgrade', 'uninstall'")

	var kataConfigResourceName string
	flag.StringVar(&kataConfigResourceName, "resource", "", "Kata Config Custom Resource Name")
	flag.Parse()

	if kataOperation == "" {
		fmt.Println("Operation type must be specified. Check -h for more information.")
		os.Exit(1)
	}
	if kataConfigResourceName == "" {
		fmt.Println("Kata Custom Resource name must be specified. Check -h for more information.")
		os.Exit(1)
	}

	switch kataOperation {
	case "install":
		installKata(kataConfigResourceName)
	case "upgrade":
		upgradeKata()
	case "uninstall":
		uninstallKata()
	default:
		fmt.Println("invalid operation. Check -h for more information.")
		os.Exit(1)
	}

}

type updateStatus = func(a *kataTypes.KataConfigStatus)

func updateKataConfigStatus(kataClientSet *kataClient.Clientset, kataConfigResourceName string, us updateStatus) (err error) {

	attempts := 5
	for i := 0; i < attempts; i++ {
		kataconfig, err := kataClientSet.KataconfigurationV1alpha1().KataConfigs(v1.NamespaceAll).Get(context.TODO(), kataConfigResourceName, metaV1.GetOptions{})
		if err != nil {
			break
		}

		us(&kataconfig.Status)

		_, err = kataClientSet.KataconfigurationV1alpha1().KataConfigs(v1.NamespaceAll).UpdateStatus(context.TODO(), kataconfig, metaV1.UpdateOptions{FieldManager: "kata-install-daemon"})
		if err != nil {
			log.Println("retrying after error:", err)
			continue
		}

		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return err
}

func getFailedNode(err error) (fn kataTypes.FailedNode, retErr error) {
	// TODO - this may not be correct. Make sure you get the right hostname
	hostname, hErr := os.Hostname()
	if hErr != nil {
		return kataTypes.FailedNode{}, hErr
	}

	return kataTypes.FailedNode{
		Name:  hostname,
		Error: fmt.Sprintf("%+v", err),
	}, nil

}

func upgradeKata() error {
	fmt.Println("Not Implemented Yet")
	return nil
}

func uninstallKata() error {
	fmt.Fprintf(os.Stderr, "%s\n", os.Getenv("PATH"))
	log.SetOutput(os.Stdout)
	cmd := exec.Command("mkdir", "-p", "/host/opt/kata-install")
	err := doCmd(cmd)
	if err != nil {
		return err
	}

	if err := syscall.Chroot("/host"); err != nil {
		log.Fatalf("Unable to chroot to %s: %s", "/host", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		log.Fatalf("Unable to chdir to %s: %s", "/", err)
	}

	cmd = exec.Command("/bin/bash", "-c", "/usr/local/kata/linux/kata-cleanup.sh")
	err = doCmd(cmd)

	return err
}

func installKata(kataConfigResourceName string) {

	//config, err := clientcmd.BuildConfigFromFlags("", "/tmp/kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		// TODO - remove printfs with logger
		fmt.Println("error creating config")
		os.Exit(1)
	}

	kataClientSet, err := kataClient.NewForConfig(config)
	if err != nil {
		fmt.Println("Unable to get client set")
		os.Exit(1)
	}
	if _, err := os.Stat("/host/usr/bin/kata-runtime"); err == nil {
		// kata exist - mark completion
		err = updateKataConfigStatus(kataClientSet, kataConfigResourceName,
			func(ks *kataTypes.KataConfigStatus) {
				if ks.InProgressNodesCount > 0 {
					ks.InProgressNodesCount = ks.InProgressNodesCount - 1
				}
				ks.CompletedNodesCount = ks.CompletedNodesCount + 1
			})

		if err != nil {
			fmt.Printf("kata exists on the node, error updating kataconfig status %+v", err)
			os.Exit(1)
		}

	} else if os.IsNotExist(err) {
		// kata doesn't exist, install it.
		err = updateKataConfigStatus(kataClientSet, kataConfigResourceName, func(ks *kataTypes.KataConfigStatus) {
			ks.InProgressNodesCount = ks.InProgressNodesCount + 1
		})

		if err != nil {
			fmt.Printf("kata is not installed on the node, error updating kataconfig status %+v", err)
			os.Exit(1)
		}

		err = installBinaries()

		// Temporary hold to simulate time taken for the installation of the binaries
		//time.Sleep(10 * time.Second)

		if err != nil {
			// kata installation failed. report it.
			err = updateKataConfigStatus(kataClientSet, kataConfigResourceName, func(ks *kataTypes.KataConfigStatus) {
				ks.InProgressNodesCount = ks.InProgressNodesCount - 1

				fn, err := getFailedNode(err)
				if err != nil {
					fmt.Printf("Error getting failed node information %+v", err)
					os.Exit(1)
				}

				ks.FailedNodes = append(ks.FailedNodes, fn)
			})

			if err != nil {
				fmt.Printf("kata installation failed, error updating kataconfig status %+v", err)
				os.Exit(1)
			}

		} else {
			// mark daemon completion
			err = updateKataConfigStatus(kataClientSet, kataConfigResourceName, func(ks *kataTypes.KataConfigStatus) {
				ks.CompletedDaemons = ks.CompletedDaemons + 1
			})

			if err != nil {
				fmt.Printf("kata installation succeeded, but error updating kataconfig status %+v", err)
				os.Exit(1)
			}
		}

	} else {
		fmt.Printf("Unknown error %+v", err)
		os.Exit(1)
	}

	for {
		c := make(chan int)
		<-c
	}
}

func doCmd(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	fmt.Println(cmd.String())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func rpmostreeOverrideReplace(rpms string) error {
	cmd := exec.Command("/bin/bash", "-c", "/usr/bin/rpm-ostree override replace /opt/kata-install/packages/"+rpms)
	err := doCmd(cmd)
	if err != nil {
		fmt.Println("override replace of " + rpms + " failed")
		return err
	}
	return nil
}

func installBinaries() error {
	fmt.Fprintf(os.Stderr, "%s\n", os.Getenv("PATH"))
	log.SetOutput(os.Stdout)
	cmd := exec.Command("mkdir", "-p", "/host/opt/kata-install")
	err := doCmd(cmd)
	if err != nil {
		return err
	}

	if err := syscall.Chroot("/host"); err != nil {
		log.Fatalf("Unable to chroot to %s: %s", "/host", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		log.Fatalf("Unable to chdir to %s: %s", "/", err)
	}

	fmt.Println("in INSTALLBINARIES")
	policy, err := signature.DefaultPolicy(nil)
	if err != nil {
		fmt.Println(err)
	}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		fmt.Println(err)
	}
	srcRef, err := alltransports.ParseImageName("docker://quay.io/jensfr/kata-artifacts:v2.0")
	if err != nil {
		fmt.Println("Invalid source name")
		return err
	}
	destRef, err := alltransports.ParseImageName("oci:/opt/kata-install/kata-image:latest")
	if err != nil {
		fmt.Println("Invalid destination name")
		return err
	}
	fmt.Println("copying down image...")
	_, err = copy.Image(context.Background(), policyContext, destRef, srcRef, &copy.Options{})
	fmt.Println("done with copying image")
	err = image.CreateRuntimeBundleLayout("/opt/kata-install/kata-image/", "/usr/local/kata", "latest", "linux", "v1.0")
	if err != nil {
		fmt.Println("error creating Runtime bundle layout in /usr/local/kata")
		return err
	}
	fmt.Println("created Runtime bundle layout in /usr/local/kata")
	fmt.Println(err)

	cmd = exec.Command("mkdir", "-p", "/etc/yum.repos.d/")
	err = doCmd(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("/usr/bin/cp", "-f", "/usr/local/kata/linux/packages.repo", "/etc/yum.repos.d/")
	err = doCmd(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("/usr/bin/cp", "-f", "/usr/local/kata/linux/katainstall.service", "/etc/systemd/system/katainstall.service")
	err = doCmd(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("/usr/bin/cp", "-f", "/usr/local/kata/linux/install_kata_packages.sh", "/opt/kata-install/install_kata_packages.sh")
	err = doCmd(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("/usr/bin/cp", "-a", "/usr/local/kata/linux/packages", "/opt/kata-install/packages")
	err = doCmd(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("env")
	err = doCmd(cmd)
	if err != nil {
		fmt.Println("calling env failed")
	}

	if err := rpmostreeOverrideReplace("linux-firmware-20191202-97.gite8a0f4c9.el8.noarch.rpm"); err != nil {
		return err
	}

	if err := rpmostreeOverrideReplace("kernel-*.rpm"); err != nil {
		return err
	}
	if err := rpmostreeOverrideReplace("{rdma-core-*.rpm,libibverbs*.rpm}"); err != nil {
		return err
	}

	cmd = exec.Command("/bin/bash", "-c", "/usr/bin/rpm-ostree install --idempotent kata-runtime kata-osbuilder")
	err = doCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}
